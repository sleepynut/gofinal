package task

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	// _ "github.com/lib/pq"
)

var DB *sql.DB

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func OpenConnection() (err error) {
	// url := "postgres://peoqxscq:o8KzOLhBc8U2tOjVkXN3g2Aj4iVSARXq@satao.db.elephantsql.com:5432/peoqxscq"
	url := os.Getenv("DATABASE_URL")
	DB, err = sql.Open("postgres", url)
	return
}

func CreateCustomer() {
	createTb := `
	CREATE TABLE IF NOT EXISTS customer (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`

	// db.Exec means to execute the given SQL
	// w/o any result|state from the result
	// following the execution
	// however the PLACEHOLDER is used to track
	// number of rows AFFECTED from the execution
	_, err := DB.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	log.Println("create table success")
}

func queryAllCustomer() []Customer {
	stmt, err := DB.Prepare("select * from customer")
	if err != nil {
		log.Fatal("can't prepare query all customer statement", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("can't query all customer", err)
	}

	var custs []Customer
	var c Customer

	for rows.Next() {
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Status); err != nil {
			log.Fatal("can't scane row into id,name,email,status", err)
			continue
		}
		custs = append(custs, c)
	}
	fmt.Println("query all customer success")
	return custs
}

func queryFilteredTodo(stat string) []Todo {
	stmt, err := DB.Prepare("select * from customer where status=$1")
	if err != nil {
		log.Fatal("can't prepare query all customer statement", err)
	}

	rows, err := stmt.Query(stat)
	if err != nil {
		log.Fatal("can't query all customer", err)
	}

	var items []Todo
	var td Todo

	for rows.Next() {
		if err := rows.Scan(&td.ID, &td.Title, &td.Status); err != nil {
			log.Fatal("can't scane row into id,title,status", err)
			continue
		}
		items = append(items, td)
	}
	fmt.Println("query filtered customer success")
	return items
}

func insertCustomer(c *Customer) *Customer {
	row := DB.QueryRow("INSERT INTO customer (name, email, status) values ($1, $2, $3)  RETURNING id;", c.Name, c.Email, c.Status)

	var id int
	err := row.Scan(&id)
	if err != nil {
		fmt.Println("can't scan id", err)
		return nil
	}
	return &Customer{id, c.Name, c.Email, c.Status}
}

func queryByID(id int) (*Customer, error) {
	var err error
	stmt, err := DB.Prepare("SELECT * FROM customer where id=$1;")
	if err != nil {
		log.Fatal("can't prepare query one row statment", err)
	}

	row := stmt.QueryRow(id)
	var c Customer

	// store results from stmt.QueryRow to each variable
	if err = row.Scan(&c.ID, &c.Name, &c.Email, &c.Status); err != nil && err != sql.ErrNoRows {
		log.Fatal("can't Scan row into variables", err)
	}

	return &c, err
}

func updateCustomer(c *Customer) {
	stmt, err := DB.Prepare("update customer set name=$2,email=$3,status=$4 where id=$1;")
	if err != nil {
		log.Fatal("can't prepare update statement", err)
	}

	if _, err := stmt.Exec(c.ID, c.Name, c.Email, c.Status); err != nil {
		log.Fatal("can't update", err)
	}
	fmt.Println("update success")
}

func deleteCustomer(id int) (err error) {
	stmt, err := DB.Prepare("delete from customer where id=$1;")
	if err != nil {
		log.Fatal("can't prepare delete statement", err)
	}

	if _, err = stmt.Exec(id); err != nil {
		log.Fatal("can't execute delete statement", err)
	}
	return
}

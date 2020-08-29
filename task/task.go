package task

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sleepynut/gofinal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Auth)

	createCustomer()

	r.GET("/customers", GetCustomersHandler)
	r.GET("/customers/:id", GetCustomersByIDHandler)

	r.POST("/customers", PostCustomersHandler)
	r.PUT("/customers/:id", UpdateCustomersHandler)

	r.DELETE("/customers/:id", DeleteCustomersHandler)
	return r
}

func GetCustomersHandler(c *gin.Context) {
	var custs []Customer
	custs = queryAllCustomer()

	c.JSON(http.StatusOK, custs)
}

func GetCustomersByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cust, err := queryByID(id)
	if err != nil {
		// DO NOT return NIL in restAPI --> not recommended
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusOK, *cust)
}

func PostCustomersHandler(c *gin.Context) {
	cust := Customer{}
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cIN, err := insertCustomer(&cust)
	if err == nil {
		c.JSON(http.StatusCreated, *cIN)
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}

}

func UpdateCustomersHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = queryByID(id)
	var cu Customer
	if err == nil {
		if err := c.ShouldBindJSON(&cu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cu.ID = id
		updateCustomer(&cu)
		c.JSON(http.StatusOK, cu)
	} else if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"warning": "Cust id NOT FOUND"})

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func DeleteCustomersHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = deleteCustomer(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

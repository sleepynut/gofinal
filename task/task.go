package task

import (
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
	cIN := insertCustomer(&cust)
	if cIN != nil {
		// b, err := json.Marshal(*cIN)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, err.Error())
		// 	return
		// }
		c.JSON(http.StatusCreated, *cIN)
	}

}

func UpdateCustomersHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cust, err := queryByID(id)
	var cu Customer
	if cust != nil {
		if err := c.ShouldBindJSON(&cu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cu.ID = id
		updateCustomer(&cu)
	}
	c.JSON(http.StatusOK, cu)
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

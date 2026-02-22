package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var permissionService = NewPermissionService()

func StartServer() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server Running"})
	})

	r.POST("/distributor", addDistributor)
	r.GET("/check/:name/:region", checkPermission)

	// 👇 NEW ENDPOINT
	r.GET("/distributor/:name", getDistributor)

	r.Run(":8080")
}

func addDistributor(c *gin.Context) {

	var distributor Distributor
	if err := c.ShouldBindJSON(&distributor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissionService.AddDistributor(&distributor)

	c.JSON(http.StatusOK, gin.H{"message": "Distributor added"})
}

func checkPermission(c *gin.Context) {

	name := c.Param("name")
	region := c.Param("region")

	allowed := permissionService.CheckPermission(name, region)

	if allowed {
		c.JSON(http.StatusOK, gin.H{"allowed": "YES"})
	} else {
		c.JSON(http.StatusOK, gin.H{"allowed": "NO"})
	}
}

func getDistributor(c *gin.Context) {

	name := c.Param("name")

	distributor, exists := permissionService.distributors[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Distributor not found",
		})
		return
	}

	c.JSON(http.StatusOK, distributor)
}

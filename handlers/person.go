package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"logtools/dbops"
	"logtools/models"
)

var db = dbops.DB

func GetPersons(c *gin.Context) {
	var users []models.Person
	if err := db.Find(&users).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, users)
	}
}

func GetPerson(c *gin.Context) {
	id := c.Param("id")
	var person models.Person
	if err := db.Where("id=?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}

func CreatePerson(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err == nil {
		db.Create(&person)
		c.JSON(200, person)
	} else {
		c.JSON(401, gin.H{"error": err, "status": "not ok"})
	}
}

func UpdatePerson(c *gin.Context) {
	var person models.Person
	id := c.Param("id")

	if err := db.Where("id=?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	fmt.Println(person)

	if err := c.ShouldBindJSON(&person); err == nil {
		db.Save(&person)
		fmt.Println(person)
		c.JSON(200, person)
	} else {
		c.JSON(401, gin.H{"error": err, "status": "not ok"})
	}
}

func DeletePerson(c *gin.Context) {
	var person models.Person
	id := c.Param("id")

	res := db.Where("id=?", id).Delete(&person)

	c.JSON(200, gin.H{"id #" + id: "deleted", "res": res})
}

package handlers

import (
	// "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"logtools/models"
)

func GetMeterConfigs(c *gin.Context) {
	var meter []models.Meter
	if err := db.Find(&meter).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, meter)
	}
}

func GetMeterConfig(c *gin.Context) {
	meterid := c.Param("meterid")
	var meter models.Meter
	if err := db.Where("meter_id=?", meterid).First(&meter).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, meter)
	}
}

func UpdateOrCreateMeterConfig(c *gin.Context) {
	meterid := c.Param("meterid")
	type Con struct {
		Config string `json:"config" binding:"required"`
	}
	var con Con
	if err := c.ShouldBindJSON(&con); err == nil {
		m := models.Meter{MeterId: meterid, Config: con.Config}
		db.Where(models.Meter{MeterId: meterid}).Assign(models.Meter{Config: con.Config}).FirstOrCreate(&m)
		c.JSON(200, gin.H{"status": "ok"})
	} else {
		c.JSON(401, gin.H{"error": err, "status": "not ok"})
	}
}

func DeleteMeterConfig(c *gin.Context) {
	var meter models.Meter
	meterid := c.Param("meterid")
	res := db.Where("meter_id=?", meterid).Delete(&meter)
	c.JSON(200, gin.H{"id #" + meterid: "deleted", "res": res, "status": "ok"})
}

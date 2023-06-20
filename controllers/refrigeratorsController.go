package controllers

import (
	"net/http"
	"refrij/initializers"
	"refrij/models"

	"github.com/gin-gonic/gin"
)

func GetRefrigeratorDetail(c *gin.Context) {
	//get id off url
	refrigerator_id := c.Param("refrigerator_id")

	//Get the posts
	var refrigerator models.Refrigerator
	result := initializers.DB.First(&refrigerator, refrigerator_id)

	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Refrigerator not found.",
		})

		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"content":      refrigerator,
		"success":      true,
		"errorMessage": nil,
	})

}

func GetRefiregators(c *gin.Context) {
	//get id off url
	user_id := c.Param("user_id")
	var refrigerators []models.Refrigerator

	result := initializers.DB.Order("created_at desc").Find(&refrigerators, "user_id = ?", user_id)
	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": result.Error.Error(),
		})
		return
	}
	if result.RowsAffected == 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "No refrigerator found",
		})
		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"content":      refrigerators,
		"success":      true,
		"errorMessage": nil,
	})

}

func CreateRefrigerator(c *gin.Context) {
	//Get the email/password of req body
	var body struct {
		RefrigeratorName string
		UserID           uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Failed to read body",
		})

		return
	}

	//Create a refrigerator
	refrigerator := models.Refrigerator{RefrigeratorName: body.RefrigeratorName, UserID: uint(body.UserID)}
	result := initializers.DB.Create(&refrigerator) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": result.Error.Error(),
		})

		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"success":      true,
		"errorMessage": nil,
	})

}
func UpdateRefrigerator(c *gin.Context) {
	//get id off url
	refrigerator_id := c.Param("refrigerator_id")

	//get the data off request body
	var body struct {
		RefrigeratorName string
	}

	c.Bind(&body)
	//Find the post were updating
	var refrigerator models.Refrigerator
	initializers.DB.First(&refrigerator, refrigerator_id)

	//Update it
	initializers.DB.Model(&refrigerator).Updates(models.Refrigerator{
		RefrigeratorName: body.RefrigeratorName,
	})
	//Respond with them
	c.JSON(200, gin.H{
		"content":      refrigerator,
		"success":      true,
		"errorMessage": nil,
	})
}

func DeleteRefrigerator(c *gin.Context) {
	//get id off url
	refrigerator_id := c.Param("refrigerator_id")

	//Delete the refrigerator
	result := initializers.DB.Unscoped().Delete(&models.Refrigerator{}, refrigerator_id)
	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": result.Error.Error(),
		})

		return
	}
	result_ingredient := initializers.DB.Unscoped().Delete(&models.Ingredient{}, "refrigerator_id = ?", refrigerator_id)

	if result_ingredient.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": result.Error.Error(),
		})

		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"success":      true,
		"errorMessage": nil,
	})
}

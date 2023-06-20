package controllers

import (
	"net/http"
	"refrij/initializers"
	"refrij/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetIngredientDetail(c *gin.Context) {
	//get id off url
	ingredient_id := c.Param("ingredient_id")

	//Get the ingredient
	var ingredient models.Ingredient
	result := initializers.DB.First(&ingredient, ingredient_id)

	if result.Error != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Record not found.",
		})

		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"content":      ingredient,
		"success":      true,
		"errorMessage": nil,
	})
}

func GetIngredients(c *gin.Context) {
	//get id off url
	refrigerator_id := c.Param("refrigerator_id")
	var ingredients []models.Ingredient

	result := initializers.DB.Order("created_at desc").Find(&ingredients, "refrigerator_id = ?", refrigerator_id)
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
			"errorMessage": "No records found",
		})
		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"content":      ingredients,
		"success":      true,
		"errorMessage": nil,
	})

}

func GetAllUserIngredients(c *gin.Context) {
	//get id off url
	user_id := c.Param("user_id")
	var ingredients []models.Ingredient

	result := initializers.DB.Unscoped().Order("created_at desc").Find(&ingredients, "user_id = ?", user_id)
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
			"errorMessage": "No records found",
		})
		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"content":      ingredients,
		"success":      true,
		"errorMessage": nil,
	})

}

func GetUserIngredients(c *gin.Context) {
	//get id off url
	user_id := c.Param("user_id")
	var ingredients []models.Ingredient

	result := initializers.DB.Unscoped().Limit(6).Order("created_at desc").Find(&ingredients, "user_id = ?", user_id)
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
			"errorMessage": "No records found",
		})
		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"content":      ingredients,
		"success":      true,
		"errorMessage": nil,
	})

}

func CreateIngredient(c *gin.Context) {
	//Get the email/password of req body
	var body struct {
		RefrigeratorID uint
		UserID         uint
		IngredientName string
		Quantity       string
		ValidUntil     time.Time `form:"end_date" binding:"required" time_format:"2006-01-02"`
	}
	c.Bind(&body)

	//Create an ingredient
	ingredient := models.Ingredient{
		RefrigeratorID: body.RefrigeratorID, UserID: body.UserID, IngredientName: body.IngredientName, Quantity: body.Quantity, ValidUntil: body.ValidUntil}

	result := initializers.DB.Create(&ingredient) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": result.Error.Error(),
		})

		return
	}

	//Return status
	c.JSON(200, gin.H{
		"success":      true,
		"errorMessage": nil,
	})

}
func UpdateIngredient(c *gin.Context) {
	//get id off url
	ingredient_id := c.Param("ingredient_id")

	//get the data off request body
	var body struct {
		IngredientName string
		Quantity       string
		ValidUntil     time.Time
	}

	c.Bind(&body)
	//Find the post were updating
	var ingredient models.Ingredient
	initializers.DB.First(&ingredient, ingredient_id)

	//Update it
	initializers.DB.Model(&ingredient).Updates(models.Ingredient{
		IngredientName: body.IngredientName,
		Quantity:       body.Quantity,
		ValidUntil:     body.ValidUntil,
	})

	//Respond with it
	c.JSON(200, gin.H{
		"content":      ingredient,
		"success":      true,
		"errorMessage": nil,
	})
}

func DeleteIngredient(c *gin.Context) {
	//get id off url
	ingredient_id := c.Param("ingredient_id")

	//Delete the refrigerator
	result := initializers.DB.Unscoped().Delete(&models.Ingredient{}, ingredient_id)

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

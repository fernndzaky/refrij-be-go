package controllers

import (
	"net/http"
	"refrij/initializers"
	"refrij/models"

	"github.com/gin-gonic/gin"
)

func GetShoppingItems(c *gin.Context) {
	//get id off url
	user_id := c.Param("user_id")
	var shoppingItems []models.Shopping

	result := initializers.DB.Unscoped().Order("created_at desc").Find(&shoppingItems, "user_id = ?", user_id)
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
		"content":      shoppingItems,
		"success":      true,
		"errorMessage": nil,
	})

}

func GetShoppingItem(c *gin.Context) {
	//get id off url
	item_id := c.Param("item_id")
	var shopping []models.Shopping

	result := initializers.DB.First(&shopping, item_id)

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
			"errorMessage": result.Error.Error(),
		})
		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"content":      shopping,
		"success":      true,
		"errorMessage": nil,
	})

}

func CreateShoppingItem(c *gin.Context) {
	//Get the email/password of req body
	var body struct {
		UserID   uint
		ItemName string
		Quantity string
		Note     string
	}
	c.Bind(&body)

	//Create a shopping item
	shoppingItem := models.Shopping{
		UserID: body.UserID, ItemName: body.ItemName, Quantity: body.Quantity, Note: body.Note}

	result := initializers.DB.Create(&shoppingItem) // pass pointer of data to Create

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

func DeleteShoppingItem(c *gin.Context) {
	//get id off url
	item_id := c.Param("item_id")

	//Delete the shopping item
	result := initializers.DB.Unscoped().Delete(&models.Shopping{}, item_id)

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

func UpdateShoppingItem(c *gin.Context) {
	//get id off url
	item_id := c.Param("item_id")

	//get the data off request body
	var body struct {
		ItemName string
		Quantity string
		Note     string
	}

	c.Bind(&body)
	//Find the post were updating
	var shoppingItem models.Shopping
	initializers.DB.First(&shoppingItem, item_id)

	//Update it
	initializers.DB.Model(&shoppingItem).Updates(models.Shopping{
		ItemName: body.ItemName,
		Quantity: body.Quantity,
		Note:     body.Note,
	})

	//Respond with it
	c.JSON(200, gin.H{
		"content":      shoppingItem,
		"success":      true,
		"errorMessage": nil,
	})

}

func UpdateIsBought(c *gin.Context) {
	//get id off url
	id := c.Param("item_id")

	//get the data off request body
	var body struct {
		IsBought bool
	}

	c.Bind(&body)
	//Find the post were updating
	var shoppingItem models.Shopping
	initializers.DB.First(&shoppingItem, id)
	//Update it
	shoppingItem.IsBought = body.IsBought
	if err := initializers.DB.Save(&shoppingItem).Error; err != nil {
		return
	}
	//Respond with it
	c.JSON(200, gin.H{
		"content":      shoppingItem,
		"success":      true,
		"errorMessage": nil,
	})
}

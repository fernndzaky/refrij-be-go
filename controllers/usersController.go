package controllers

import (
	"errors"
	"net/http"
	"os"
	"refrij/initializers"
	"refrij/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"encoding/hex"
	"math/rand"
)

const (
	randomStringLength = 6
)

func generateRandomID() (string, error) {
	for {
		// Generate a random 6-byte string
		bytes := make([]byte, 3)
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}

		// Convert the bytes to a hex string
		id := hex.EncodeToString(bytes)

		// Check if the generated ID already exists
		exists, err := checkIDExists(id)
		if err != nil {
			return "", err
		}

		// If the ID is unique, return it
		if !exists {
			return id, nil
		}
	}
}

func checkIDExists(id string) (bool, error) {
	// Use your existing logic to check if the ID exists
	var user models.User
	err := initializers.DB.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Login(c *gin.Context) {
	//Get the email and paass off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Failed to read body..",
		})

		return
	}

	// Lookk up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Email not found.",
		})

		return
	}

	// Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Ouch.. Wrong password.",
		})

		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Failed to create token",
		})

		return
	}
	// send it back

	// Set the token as Bearer token in the response header
	c.Header("Authorization", "Bearer "+tokenString)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"errorMessage": nil,
		"token":        tokenString,
		"content":      user,
	})

}

func SignUp(c *gin.Context) {
	//Allow CORS here By * or specific origin
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
	//Get the email/password of req body
	var body struct {
		Email    string
		Name     string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Failed to read body..",
		})

		return
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Failed to hash the password..",
		})

		return
	}

	//Create a user
	user := models.User{Email: body.Email, Name: body.Name, Password: string(hash)}
	result := initializers.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": result.Error.Error(),
		})

		return
	}

	//Return status
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"errorMessage": nil,
	})

}

func GetUserDetail(c *gin.Context) {

	//get id off url
	user_id := c.Param("user_id")

	//Get the posts
	var user models.User
	result := initializers.DB.First(&user, "id = ?", user_id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": result.Error.Error(),
		})

		return
	}

	//Respond with them
	c.JSON(200, gin.H{
		"content":      user,
		"success":      true,
		"errorMessage": nil,
	})
}

func UpdateProfile(c *gin.Context) {
	//get id off url
	user_id := c.Param("user_id")

	//get the data off request body
	var body struct {
		Name string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Failed to read body..",
		})

		return
	}

	// Validate the length of the Name field
	if len(body.Name) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Name must not exceed 50 characters.",
		})
		return
	}

	//Find the user were updating
	var user models.User
	initializers.DB.First(&user, "id = ?", user_id)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Email not found.",
		})

		return
	}

	//Update it
	result := initializers.DB.Model(&user).Updates(models.User{
		Name: body.Name,
	})

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": result.Error.Error(),
		})

		return
	}

	//Respond with it
	c.JSON(200, gin.H{
		"content":      user,
		"success":      true,
		"errorMessage": nil,
	})

}

func ChangePassword(c *gin.Context) {
	//get id off url
	id := c.Param("id")

	//get the data off request body
	var body struct {
		OldPassword        string
		NewPassword        string
		ConfirmNewPassword string
	}

	c.Bind(&body)
	//Find the post were updating
	var user models.User
	initializers.DB.First(&user, "id = ?", id)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "User not found..",
		})

		return
	}

	// Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Ouch.. Wrong old password",
		})

		return
	}

	if body.NewPassword != body.ConfirmNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "New password dont match",
		})

		return
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": "Failed to hash the password",
		})

		return
	}

	//Update it
	initializers.DB.Model(&user).Updates(models.User{
		Password: string(hash),
	})

	//Respond with it
	c.JSON(200, gin.H{
		"success":      true,
		"errorMessage": nil,
	})
}

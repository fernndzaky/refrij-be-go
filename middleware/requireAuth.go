package middleware

import (
	"fmt"
	"net/http"
	"os"
	"refrij/initializers"
	"refrij/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("In Middleware")

	// Get the Authorization header from the request
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract the token from the Authorization header
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenString := tokenParts[1]

	// Decode / validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g., []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fmt.Println(claims["sub"])

		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user, "id = ?", claims["sub"])

		if user.ID == 0 {
			fmt.Println("error here 2")
			fmt.Println(user)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Attach to the request context
		c.Set("user", user)

		// Continue
		c.Next()
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

package initializers

import (
	"fmt"
	"refrij/models"
)

func SyncDatabase() {
	fmt.Println("migrating")
	DB.AutoMigrate(&models.User{}, &models.Shopping{}, &models.Refrigerator{}, &models.Ingredient{}, &models.Category{})
}

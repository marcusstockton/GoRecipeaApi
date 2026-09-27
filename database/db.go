package database

import (
	"errors"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"recipea.com/m/domain/chef"
	"recipea.com/m/domain/recipe"
)

var DB *gorm.DB

func defaultDSN() string {
	return "host=localhost user=recipea password=recipea dbname=recipea_dev port=5432 sslmode=disable"
}

func InitDB(dsn string) {
	if dsn == "" {
		dsn = defaultDSN()
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(gormLogModeFromAppLogLevel(os.Getenv("LOG_LEVEL"))),
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}
}

func gormLogModeFromAppLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "warn", "warning":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Warn
	}
}

func Seed() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var chefCount int64
		if err := tx.Model(&chef.Chef{}).Count(&chefCount).Error; err != nil {
			return err
		}
		if chefCount == 0 {
			seededChef, err := chef.NewChef("Gordon", "Ramsay", "gordon@example.com", "securePassword123")
			if err != nil {
				return err
			}
			if err := tx.Create(seededChef).Error; err != nil {
				return err
			}
		}

		var recipeCount int64
		if err := tx.Model(&recipe.Recipe{}).Count(&recipeCount).Error; err != nil {
			return err
		}
		if recipeCount > 0 {
			return nil
		}

		var seededChef chef.Chef
		if err := tx.Order("id asc").First(&seededChef).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		seededRecipe := recipe.Recipe{
			Title:       "Chocolate Cake",
			Description: "A delicious homemade chocolate cake",
			OwnedBy:     seededChef.ID,
			Ingredients: []recipe.Ingredient{
				{Quantity: "2", Unit: "cups", Name: "flour", Preparation: "sifted"},
				{Quantity: "1", Unit: "cup", Name: "sugar"},
			},
			Steps: []recipe.Step{
				{Order: 1, Action: "Preheat oven to 350°F"},
				{Order: 2, Action: "Mix dry ingredients in a bowl"},
			},
		}

		return tx.Create(&seededRecipe).Error
	})
}

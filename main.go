package main

import (
	"Gofiber/models"
	"Gofiber/service"
	"Gofiber/storage"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main () {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate database")
	}
	r := service.Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetUpRoutes(app)
	app.Get("/", func(context *fiber.Ctx) error {
		return context.SendString("Welcome to Gofiber")
	})

	if os.Getenv("PORT") == "" {
		app.Listen(":3000")
	} else {
		app.Listen(":"+os.Getenv("PORT"))
	}
}

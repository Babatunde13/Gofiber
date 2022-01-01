package service

import (
	"Gofiber/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Books struct {
	Author   string `json:"author" validate:"required"`
	Publisher string `json:"publisher" validate:"required"`
	Title    string `json:"title" validate:"required"`
}

type Repository struct {
	DB *gorm.DB
}


func (r *Repository) createBook(context *fiber.Ctx) error {
	book := Books{}
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid book data",
			"status": false,
		})
		return nil
	}
	// validate book data
	validator := validator.New()
	err = validator.Struct(book)
	if err != nil {
		context.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid book data",
			"status": false,
		})
		return nil
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error creating book",
			"status": false,
		})
		return nil
	}
	context.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Book created successfully",
		"status": true,
		"data": book,
	})
	return nil
}

func (r *Repository) updateBookById(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		context.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid book data",
			"status": false,
		})
		return nil
	}
	bookModel := &models.Book{}
	book := Books{}
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid book data",
			"status": false,
		})
		return nil
	}
	err = r.DB.First(bookModel, id).Error
	if err != nil {
		context.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "Book not found",
			"status": false,
		})
		return nil
	}
	if book.Author != "" {
		bookModel.Author = book.Author
	}
	if book.Publisher != "" {
		bookModel.Publisher = book.Publisher
	}
	if book.Title != "" {
		bookModel.Title = book.Title
	}
	err = r.DB.Save(bookModel).Error
	if err != nil {
		context.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Error updating book",
			"status": false,
		})
		return nil
	}
	context.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Book updated successfully",
		"status": true,
		"data": bookModel,
	})
	return nil
}


func (r *Repository) getBookById(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		context.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid book data",
			"status": false,
		})
		return nil
	}
	bookModel := &models.Book{}
	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "Book not found",
			"status": false,
		})
		return nil
	}
	context.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Book fetched successfully",
		"status": true,
		"data": bookModel,
	})
	return nil
}

func (r *Repository) getAllBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Book{}
	err := r.DB.Find(&bookModels).Error
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error fetching books",
			"status": false,
		})
	}
	context.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Books fetched successfully",
		"status": true,
		"data": bookModels,
	})
	return nil
}

func (r *Repository) deleteBookById(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		context.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "Invalid book data",
			"status": false,
		})
		return nil
	}
	bookModel := &models.Book{}
	err := r.DB.Delete(&bookModel, id).Error
	if err != nil {
		context.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error deleting book",
			"status": false,
		})
		return nil
	}
	context.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Book deleted successfully",
		"status": true,
		"data": nil,
	})
	return nil
}

func (r *Repository) SetUpRoutes (app *fiber.App) {
	api := app.Group("/api/v1")
	api.Get("/docs", func(context *fiber.Ctx) error {
		return context.Status(fiber.StatusMovedPermanently).Redirect("https://documenter.getpostman.com/view/11853513/UVRHi3Jn")
	})
	api.Post("/books", r.createBook)
	api.Get("/books", r.getAllBooks)
	api.Get("/books/:id", r.getBookById)
	api.Put("/books/:id", r.updateBookById)
	api.Delete("/books/:id", r.deleteBookById)
}
package routes

import (
	"strconv"
	"todo-app/app/middleware"
	"todo-app/config/database"

	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	return middleware.Redirect(c, "index", "/")
}

func LoginPost(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if !database.UserExists(email, "email") {
		return c.SendString("Wrong email")
	}

	user := database.SearchUserByString(email, "email")

	// Check if the password matches the password hash
	if middleware.ValidatePassword(user.Password, c.FormValue("password")) {
		middleware.SetSessionCookie(c, user.ID)
		return middleware.Redirect(c, "index", "/")
	}

	return c.SendString("Wrong password")
}

func LoginGet(c *fiber.Ctx) error {
	return middleware.Redirect(c, "login", "/login")
}

func RegisterPost(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if !database.UserExists(email, "email") {
		user := database.User{
			Username: c.FormValue("username"),
			Email:    c.FormValue("email"),
			Password: middleware.HashPassword(c.FormValue("password")),
		}

		database.Database.Create(&user)
		return middleware.Redirect(c, "login", "/login")
	}

	return c.SendString("Email already used")
}

func RegisterGet(c *fiber.Ctx) error {
	return middleware.Redirect(c, "register", "/register")
}

func LogoutGet(c *fiber.Ctx) error {
	middleware.ClearSessionCookie(c)
	return middleware.Redirect(c, "index", "/")
}

func TodoGet(c *fiber.Ctx) error {
	return middleware.Redirect(c, "todo", "/todo")
}

func ListAddPost(c *fiber.Ctx) error {
	list := database.List{
		Title:  c.FormValue("list_title"),
		UserID: c.Locals("user_id").(uint),
	}

	database.Database.Create(&list)
	// Quando for uma requisição HTMX, use RenderPartial
	return middleware.Render(c, "partials/menus-list", true)
}

func ListDeletePost(c *fiber.Ctx) error {
	// Parse the list ID from the request
	listID, err := strconv.ParseUint(c.FormValue("list_id"), 10, 32)
	if err != nil {
		return c.Status(400).SendString("Invalid list ID")
	}

	// Get the user ID from the context
	userID := c.Locals("user_id").(uint)

	// Find the list
	var list database.List
	if err := database.Database.First(&list, listID).Error; err != nil {
		return c.Status(404).SendString("List not found")
	}

	// Verify that the list belongs to the user
	if list.UserID != userID {
		return c.Status(403).SendString("Unauthorized")
	}

	// Delete all tasks associated with the list (GORM will handle this automatically due to the relationship)
	if err := database.Database.Delete(&list).Error; err != nil {
		return c.Status(500).SendString("Failed to delete list")
	}

	// Return the updated lists partial
	return middleware.Render(c, "partials/menus-list", true)
}

func TaskAddPost(c *fiber.Ctx) error {
	listID, err := strconv.ParseUint(c.FormValue("list_id"), 10, 32)
	if err != nil {
		return c.Status(400).SendString("Invalid list ID")
	}

	task := database.Task{
		Title:       c.FormValue("task_title"),
		Description: c.FormValue("task_description"),
		DueDate:     c.FormValue("task_date"),
		Completed:   false,
		ListID:      uint(listID),
	}

	// todo, task e lista direito, tem coisas erradas

	database.Database.Create(&task)
	return c.SendString("Task created")
}

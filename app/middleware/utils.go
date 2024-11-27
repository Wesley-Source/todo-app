package middleware

import (
	"log"
	"os"
	"time"
	"todo-app/config/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3/v2"
	"github.com/gofiber/utils"
	"golang.org/x/crypto/bcrypt"
)

var Session *session.Store

type List struct {
	ID    uint
	Title string
	Tasks []Task
}

type Task struct {
	ID          uint
	Title       string
	Description string
	DueDate     string
	Completed   bool
}

func listConverter(list database.List) List {
	return List{
		ID:    list.ID,
		Title: list.Title,
		Tasks: []Task{},
	}
}

func taskConverter(task database.Task) Task {
	return Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Completed:   task.Completed,
	}
}

func listMaker(user database.User) []List {
	var lists []List
	for _, list := range database.SearchListsByUserID(user.ID) {
		convertedList := listConverter(list)
		for _, task := range database.SearchTasksByListID(list.ID) {
			convertedTask := taskConverter(task)
			convertedList.Tasks = append(convertedList.Tasks, convertedTask)
		}

		lists = append(lists, convertedList)
	}

	return lists
}

func Render(c *fiber.Ctx, view string, partial ...bool) error {
	userID := GetSessionCookie(c)
	data := fiber.Map{
		"Title":  os.Getenv("TITLE"),
		"UserID": userID,
	}

	if userID != nil {
		// If the user is logged in, fetch their information
		user := database.SearchUserById(userID.(uint))

		lists := listMaker(user)

		data["Username"] = user.Username
		data["Email"] = user.Email
		data["Lists"] = lists
	}

	if len(partial) > 0 && partial[0] {
		return c.Render(view, data)
	}

	return c.Render(view, data, "layouts/main")
}

func Redirect(c *fiber.Ctx, view, route string) error {
	if c.Get("HX-Request") == "true" {
		c.Set("HX-Redirect", route)
		return c.SendStatus(fiber.StatusOK)
	}

	return Render(c, view)
}

func ConnectSessionsDB() {
	storage := sqlite3.New(sqlite3.Config{
		Table:    "session_storage",
		Database: "./config/database/sessions.db",
	})
	Session = session.New(session.Config{
		Storage:        storage,
		Expiration:     24 * time.Hour,
		KeyLookup:      "cookie:session_id",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "Strict",
		KeyGenerator:   utils.UUID,
	})
}

func HashPassword(password string) string {
	// Returns a hashed and salted password
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashBytes)
}

func ValidatePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func AuthMiddleware(c *fiber.Ctx) error {
	/*
		Prevents the logged user from accessing the login and register pages
		and prevents non-logged users from accessing protected pages
	*/
	sess, err := Session.Get(c)
	if err != nil {
		log.Println(err)
	}

	userID := sess.Get("user_id")
	switch c.Path() {

	// If the user is logged in, redirect them to the index page instead of /login or /register
	case "/login", "/register":

		// Redirects the user to the index page if they're already logged in
		if userID != nil {
			return Redirect(c, "index", "/")
		} else {
			return c.Next()
		}
	}

	// If the user is not logged in
	if userID == nil {
		return Redirect(c, "index", "/")
	}

	sess.Save()
	c.Locals("user_id", userID.(uint))

	return c.Next()
}

func SetSessionCookie(c *fiber.Ctx, id uint) {
	session, err := Session.Get(c)
	if err != nil {
		log.Println("Failed to get session.")
	}

	// Saves the user_id as a cookie in the user's browser
	session.Set("user_id", id)
	session.Save()
}

func GetSessionCookie(c *fiber.Ctx) interface{} {
	session, err := Session.Get(c)
	if err != nil {
		log.Println("Failed to get session.")
	}

	return session.Get("user_id")
}

func ClearSessionCookie(c *fiber.Ctx) {
	session, err := Session.Get(c)
	if err != nil {
		log.Println("Failed to get session.")
	} else {
		session.Delete("user_id")
		session.Save()
	}
	c.ClearCookie("user_id")
}

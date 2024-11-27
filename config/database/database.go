package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Models:
type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Lists    []List
}

type List struct {
	gorm.Model
	Title  string `json:"title"`
	UserID uint
	Tasks  []Task
}

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Completed   bool   `json:"completed" gorm:"default:false"`
	ListID      uint
}

var Database *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("./config/database/models.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&User{}, &List{}, &Task{}) // Include models as you create them

	Database = db // Moving the variable to the global scope
}

func UserExists(value, condition string) bool {
	var user User

	// Checks if user exists by email
	conditionString := condition + " = ?"
	result := Database.Where(conditionString, value).First(&user)

	return result.Error != gorm.ErrRecordNotFound // Returns true if finds a user, false if don't
}

func SearchUserByString(value, condition string) User {
	// Search user by any string type value

	var user User

	conditionString := condition + " = ?"
	Database.Where(conditionString, value).First(&user)
	return user
}

func SearchUserById(id uint) User {
	var user User

	Database.Where("id = ?", id).First(&user)
	return user
}

func SearchListsByUserID(userID uint) []List {
	// Search lists by user's ID

	var lists []List

	Database.Where("user_id = ?", userID).Find(&lists)
	return lists
}

func SearchTasksByListID(listID uint) []Task {
	// Search tasks by list's ID

	var tasks []Task

	Database.Where("list_id = ?", listID).Find(&tasks)
	return tasks
}

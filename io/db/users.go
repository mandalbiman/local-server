package db

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"bats.com/local-server/api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserDB struct {
	db *gorm.DB
}

func GetUserDB() *UserDB {
	db := setUpDB()
	db.AutoMigrate(&models.User{}) // ensure table exists
	return &UserDB{db: db}
}

func setUpDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

// Get all users
func (u UserDB) GetUsers() ([]models.User, error) {
	var users []models.User
	result := u.db.Find(&users)
	return users, result.Error
}

// Get user by ID
func (u UserDB) GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := u.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create new user
func (u UserDB) CreateUser(user models.User) (*models.User, error) {
	result := u.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Update existing user
func (u UserDB) UpdateUser(user models.User) (*models.User, error) {
	result := u.db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Delete user
func (u UserDB) DeleteUser(id int) (*models.User, error) {
	user, err := u.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	result := u.db.Delete(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// Search by username
func (u UserDB) SearchUserByUsername(username string) ([]models.User, error) {
	var users []models.User
	result := u.db.Where("username LIKE ?", "%"+username+"%").Find(&users)
	return users, result.Error
}

func (u UserDB) LoadUsersFromJSON(file string) error {
	// Read JSON file
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Temporary struct to handle CreatedAt as string
	type UserJSON struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Age       int    `json:"age"`
		Gender    string `json:"gender"`
		Phone     string `json:"phone"`
		City      string `json:"city"`
		Country   string `json:"country"`
		CreatedAt string `json:"created_at"`
	}

	var usersJSON []UserJSON
	if err := json.Unmarshal(data, &usersJSON); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to models.User
	var users []models.User
	for _, ujson := range usersJSON {
		// Parse timestamp without timezone
		t, err := time.Parse("2006-01-02T15:04:05.999999", ujson.CreatedAt)
		if err != nil {
			t = time.Now() // fallback if parse fails
		}

		users = append(users, models.User{
			Name:      ujson.Name,
			Username:  ujson.Username,
			Email:     ujson.Email,
			Password:  ujson.Password,
			Age:       ujson.Age,
			Gender:    ujson.Gender,
			Phone:     ujson.Phone,
			City:      ujson.City,
			Country:   ujson.Country,
			CreatedAt: t,
		})
	}

	// Insert users into DB
	for i := range users {
		_, err := u.CreateUser(users[i])
		if err != nil {
			fmt.Printf("Failed to insert user %s: %v\n", users[i].Username, err)
		}
	}

	fmt.Printf("✅ Loaded %d users into SQLite\n", len(users))
	return nil
}

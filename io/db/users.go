package db

import (
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
func (u UserDB) SearchUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := u.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

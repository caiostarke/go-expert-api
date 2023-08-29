package database

import (
	"testing"

	"github.com/caiostarke/go-expert-apo/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("John", "j@j.com", "123456")
	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Email, userFound.Email)
	assert.Equal(t, user.Name, userFound.Name)
	assert.NotNil(t, userFound.Password)
}

func TestGetUserByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("John", "j@j.com", "123456")

	userDB := NewUser(db)
	err = userDB.Create(user)
	assert.Nil(t, err)

	userFound, err := userDB.FindByEmail("j@j.com")
	assert.Nil(t, err)
	assert.NotNil(t, userFound)

	assert.Equal(t, userFound.Email, user.Email)
	assert.Equal(t, userFound.ID, user.ID)
	assert.Equal(t, userFound.Name, user.Name)
	assert.NotEmpty(t, userFound.Password)
}

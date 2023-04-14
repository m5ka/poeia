package poeia

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
	Hash     string
}

type UserCreate struct {
	Username string `validate:"min=3,max=32"`
	Email    string `validate:"email"`
	Password string `validate:"min=8"`
}

func (db *DB) FindUser(ctx context.Context, username string) (user *User, err error) {
	err = db.db.WithContext(ctx).Find(&user, "username = ?", username).Error
	return
}

func (db *DB) CreateUser(ctx context.Context, create UserCreate) (user *User, err error) {
	user = &User{}
	return
}

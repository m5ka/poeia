package poeia

import (
	"fmt"

	"github.com/go-playground/validator"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	db       *gorm.DB
	validate *validator.Validate

	Host     string
	User     string
	Password string
	Database string
}

func Open(host string, user string, password string, database string) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})

	return &DB{
		db:       db,
		validate: validator.New(),

		Host:     host,
		User:     user,
		Password: password,
		Database: database,
	}, nil
}

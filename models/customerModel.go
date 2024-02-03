package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	FullName    string
	Address     string
	PhoneNumber string `gorm:"unique"`
	UserID      uint32
	User        User `gorm:"references:id"`
}

package models

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	NoSiup string `grom:"unique"`
	Name string
	Address string
	PhoneNumber string `grom:"unique"`
}
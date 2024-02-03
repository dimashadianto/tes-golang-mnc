package models

import "gorm.io/gorm"

type Seller struct {
	gorm.Model
	FullName    string
	Email       string `grom:"unique"`
	PhoneNumber string `grom:"unique"`
	Address     string
	StoreID     uint32
	Store       Store `gorm:"references:id"`
}

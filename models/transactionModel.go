package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	CustomerID             uint32
	Customer               Customer `gorm:"references:id"`
	TransactionAmount      uint64
	TransactionDescription string
	SellerID               uint32
	Seller                 Seller `gorm:"references:id"`
}

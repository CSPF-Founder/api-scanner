package models

import "gorm.io/gorm"

var db *gorm.DB

func SetupDB(d *gorm.DB) {
	db = d
}

// Response contains the attributes found in an API response
type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

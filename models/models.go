package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string
	StudentID uint
	Email     string `gorm:"unique"`
	Password  string
	Role      string
}

type Leave struct {
	gorm.Model
	StudentID uint
	StartDate string
	EndDate   string
	Reason    string
	Status    string
}

type Attendance struct {
	gorm.Model
	StudentID uint
	Date      string
	Present   bool
}

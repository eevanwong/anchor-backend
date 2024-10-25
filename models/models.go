package models

import (
	"gorm.io/gorm"
)

/*
// gorm.Model definition -> built in
type Model struct {
  ID        uint           `gorm:"primaryKey"` <- auto increments from 1
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}
*/
// Define your models
type User struct {
	gorm.Model
	Name        string `gorm:"<-"`
	Contact     string `gorm:"<-"`
	ContactType string `gorm:"<-"`
}

type Rack struct {
	gorm.Model
	CurrUserID uint `gorm:"<-"`
}

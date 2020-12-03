package model

import (
	"time"

	"gorm.io/datatypes"
)

type Company struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:""`
	// Track creating/updating time/unix (milli/nano) seconds for multiple fields
	CreatedAt time.Time // Set to current time if it is zero on creating
	UpdatedAt int       // Set to current unix seconds on updating or if it is zero on creating
	Updated   int64     `gorm:"autoUpdateTime:nano"`  // Use unix Nano seconds as updating time
	Updated2  int64     `gorm:"autoUpdateTime:milli"` // Use unix Milli seconds as updating time
	Created   int64     `gorm:"autoCreateTime"`       // Use unix seconds as creating time
}

type User struct {
	ID         int `gorm:"primaryKey"`
	CompanyID  int
	Company    Company
	Name       string
	Address    string
	Age        int
	CreditCard CreditCard
	Posts      []Post
	Attributes datatypes.JSON
}

type CreditCard struct {
	ID     int `gorm:"primaryKey"`
	UserID int
	Number string
}

type Post struct {
	ID      int `gorm:"primaryKey"`
	UserID  int
	Content string
	Tags    []Tag `gorm:"many2many:post_tags"`
}

type Tag struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

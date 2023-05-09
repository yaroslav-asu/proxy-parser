package models

type Site struct {
	ID  uint `gorm:"primaryKey"`
	Url string
}

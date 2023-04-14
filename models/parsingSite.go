package models

type Site struct {
	ID  uint `gorm:"primaryKey"`
	Url string
}

func NewSite(url string) Site {
	return Site{
		Url: url,
	}
}

package models

type AvailablePlans struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float32
	Recurrence  int
}

package models

import "time"

type AvailablePlans struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float32
	Recurrence  int
}

type UserPlans struct {
	Id             string `json:"id"`
	UserId         string `json:"userid"`
	PlanId         int
	Status         string
	CreatedTs      time.Time
	LastModifiedTs time.Time
}

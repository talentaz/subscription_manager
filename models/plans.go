package models

import "time"

type AvailablePlans struct {
	Id          int `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float32
	Recurrence  int
	PriceId     string `json:"priceid"`
}

type UserPlans struct {
	Id             string `json:"id"`
	UserId         string `json:"userid"`
	PlanId         int
	CustomerId     string `json:"customerid"`
	PriceId        string `json:"priceid"`
	SubscriptionId string `json:"subscriptionid"`
	Status         string
	CreatedTs      time.Time
	LastModifiedTs time.Time
}

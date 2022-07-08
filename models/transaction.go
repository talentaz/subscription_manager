package models

import "time"

type Transactions struct {
	Id             string `json:"id"`
	UserId         string `json:"userid"`
	PriceId        string `json:"priceid"`
	UserPlanId     string `json:"userpalnid"`
	SessionId      string `json:"sessionid"`
	CustomerId     string `json:"customerid"`
	Status         string
	CreatedTs      time.Time
	LastModifiedTs time.Time
}

package gw2models

import "time"

type ActiveTransaction struct {
	Id       int       `json:"id"`
	ItemId   int       `json:"item_id"`
	Price    int       `json:"price"`
	Quantity int       `json:"quantity"`
	Created  time.Time `json:"created"`
}

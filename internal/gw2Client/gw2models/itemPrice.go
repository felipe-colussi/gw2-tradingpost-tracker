package gw2models

type Order struct {
	Quantity  int `json:"quantity"`
	UnitPrice int `json:"unit_price"`
}

type OrderObject struct {
	Id          int  `json:"id"`
	Whitelisted bool `json:"whitelisted"`
	Buys        Order
	Sells       Order
}

package runes

type charmPrice struct {
	Sell    int
	Buy     int
	SellTax int
	BuyTax  int
}

type runeToBuy struct {
	Name         string
	Id           int
	InstaBuyRune int
	OrderBuyRune int
	InstaBuy     Proffit
	OrderBuy     Proffit
}

type Proffit struct {
	InstaSell      int
	InstaSellTpTax int
	OrderSell      int
	OrderSellTpTax int
}

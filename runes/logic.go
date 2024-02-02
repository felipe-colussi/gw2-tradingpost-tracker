package runes

import (
	"errors"
	"fmt"
	"gw2-tradingpost-tracker/internal/gw2Client"
	"gw2-tradingpost-tracker/internal/gw2Client/gw2models"
	"sort"
	"sync"
)

type RuneSerarchInterface interface {
	ConvertRunesIntoCharmsAndSell(searchSkill, searchPortence, searchBrilliance bool)
}

func New(apiKey string) (RuneSerarchInterface, error) {
	if apiKey == "" {
		return nil, errors.New("API KEY NEEDED")
	}
	return &runeSearcher{
		c:          &gw2Client.Client{Apikey: apiKey},
		charmPrice: map[int]charmPrice{},
	}, nil
}

type runeSearcher struct {
	c          *gw2Client.Client
	charmPrice map[int]charmPrice
}

var _ RuneSerarchInterface = &runeSearcher{}

func (r *runeSearcher) ConvertRunesIntoCharmsAndSell(searchSkill, searchPortence, searchBrilliance bool) {

	err := r.getCharmPrice()
	if err != nil {
		fmt.Println(err)
		return
	}

	transactions, err := r.c.BuyTransactions()
	transactionMap := transactionsToMap(transactions)

	wg := sync.WaitGroup{}

	if searchPortence {
		wg.Add(1)
		go func() {
			r.searchFor(potenceList, r.charmPrice[CharmOfPotence], transactionMap, "potence")
			wg.Done()
		}()

	}

	if searchBrilliance {
		wg.Add(1)
		go func() {
			r.searchFor(brillianceList, r.charmPrice[CharmOfBrilliance], transactionMap, "brilliance")
			wg.Done()
		}()

	}

	if searchSkill {
		wg.Add(1)
		go func() {
			r.searchFor(skillIst, r.charmPrice[CharmOfSkill], transactionMap, "skill")
			wg.Done()
		}()

	}
	wg.Wait()

}

func (r *runeSearcher) searchFor(arrayOfRunes []int, charmPrice charmPrice, transactionMap map[int]int, name string) {
	runesTobuy, err := r.getWorthRunes(charmPrice, arrayOfRunes)

	if err != nil {
		fmt.Printf("Failed to validate name for rune %s, with error %s\n", name, err.Error())
		return
	}

	if len(runesTobuy) == 0 {
		fmt.Println("=======================================================")
		fmt.Println(name)
		fmt.Println("No rune worth buying to craft this one.")
		return
	}

	if err := r.populateNames(runesTobuy); err != nil {
		fmt.Printf("Failed to populate names for rune %s, with error %s\n", name, err.Error())
		return
	}
	PrintArray(name, runesTobuy, transactionMap)
}

func (r *runeSearcher) populateNames(toBuy []runeToBuy) error {
	m := map[int]string{}

	itemsId := []int{}
	for _, v := range toBuy {
		itemsId = append(itemsId, v.Id)
	}

	items, err := r.c.GetItens(itemsId)
	if err != nil {
		return err
	}
	for _, v := range items {
		m[v.GetId()] = v.GetName()
	}
	for k, v := range toBuy {
		toBuy[k].Name = m[v.Id]
	}
	return nil
}

func (r *runeSearcher) getWorthRunes(charm charmPrice, savagables []int) ([]runeToBuy, error) {
	runes, err := r.c.GetPriceOfMultiple(savagables)
	if err != nil {
		return nil, err
	}
	buy := []runeToBuy{}

	for _, superiorRune := range runes {
		if superiorRune.Buys.UnitPrice*10 < int(float64(charm.Sell)*0.85) {
			buy = append(buy, calcProffit(charm, superiorRune))
		}
	}
	return buy, nil
}

func (r *runeSearcher) getCharmPrice() error {
	if r.charmPrice == nil {
		r.charmPrice = map[int]charmPrice{}
	}

	charms, err := r.c.GetPriceOfMultiple([]int{CharmOfBrilliance, CharmOfPotence, CharmOfSkill})
	if err != nil {
		return err
	}

	for _, c := range charms {
		r.charmPrice[c.Id] = charmPrice{
			Sell:    c.Sells.UnitPrice,
			Buy:     c.Buys.UnitPrice,
			SellTax: func() int { return int(float64(c.Sells.UnitPrice) * .85) }(),
			BuyTax:  func() int { return int(float64(c.Buys.UnitPrice) * .85) }(),
		}
	}
	return nil
}

func calcProffit(c charmPrice, r gw2models.OrderObject) runeToBuy {
	buy := runeToBuy{
		Id:           r.Id,
		InstaBuyRune: r.Sells.UnitPrice,
		OrderBuyRune: r.Buys.UnitPrice + 1, // add on gold to cut oreder
	}
	buy.InstaBuy = Proffit{ // Buying Using sell price, insta buy
		InstaSell:      c.Buy - (r.Sells.UnitPrice * 10), // Sell From Buy price (instasell) less buy for sell (insta buy)
		InstaSellTpTax: int(float64(c.Buy-(r.Sells.UnitPrice*10)) * 0.85),
		OrderSell:      c.Sell - (r.Sells.UnitPrice * 10),
		OrderSellTpTax: int(float64(c.Sell-(r.Sells.UnitPrice*10)) * 0.85),
	}

	buy.OrderBuy = Proffit{
		InstaSell:      c.Buy - (r.Buys.UnitPrice * 10), // Sell From Buy price (instasell) less buy for Buy (order)
		InstaSellTpTax: int(float64(c.Buy-(r.Buys.UnitPrice*10)) * 0.85),
		OrderSell:      c.Sell - (r.Buys.UnitPrice * 10),
		OrderSellTpTax: int(float64(c.Sell-(r.Buys.UnitPrice*10)) * 0.85),
	}
	return buy
}

func PrintArray(t string, a []runeToBuy, myTransactions map[int]int) {
	sort.Slice(a, func(i, j int) bool {
		return a[i].OrderBuy.OrderSellTpTax > a[j].OrderBuy.OrderSellTpTax
	})
	fmt.Println("==========================================================================")
	fmt.Println(t)
	fmt.Println("\n")

	for _, v := range a {

		defaultCollor := "\033[37m"
		if v.OrderBuyRune-1 == myTransactions[v.Id] {
			fmt.Print(printRed)
			defaultCollor = printRed
		}

		fmt.Printf("%s - InstaBuy: %d, Order: %s %d %s, MyBuyOrder:  %s %d %s Proffits: Order_buy_sel: %s %d %s,  order_buy_insta_sell: %d, insta_buy_order_sell: %d  \n\n",
			v.Name, v.InstaBuyRune, printYellow, v.OrderBuyRune, defaultCollor, printYellow, myTransactions[v.Id], defaultCollor, printGreen, v.OrderBuy.OrderSellTpTax/10, defaultCollor, v.OrderBuy.InstaSellTpTax/10, v.InstaBuy.OrderSellTpTax/10)

		fmt.Print(printReset)

	}

}

func transactionsToMap(t []gw2models.ActiveTransaction) map[int]int {
	m := make(map[int]int, len(t))
	for _, v := range t {
		if m[v.ItemId] < v.Price {
			m[v.ItemId] = v.Price
		}
	}
	return m
}

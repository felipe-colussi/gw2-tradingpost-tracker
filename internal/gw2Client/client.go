package gw2Client

import (
	"encoding/json"
	"errors"
	"fmt"
	"gw2-tradingpost-tracker/internal/gw2Client/gw2models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	Apikey string
}

func (c Client) ListAllIds() ([]int, error) {
	r, err := http.Get("https://api.guildwars2.com/v2/items")
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.New("Non 200 status code")
	}
	ids := []int{}
	responseBody, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		return nil, errReadAll
	}

	if err := json.Unmarshal(responseBody, &ids); err != nil {
		return nil, err
	}

	return ids, nil

}

func (c Client) GetItens(ids []int) ([]gw2models.BaseItemData, error) {

	if ids == nil || len(ids) == 0 {
		return nil, errors.New("No Ids passed")
	}
	urlBuilder := strings.Builder{}
	urlBuilder.WriteString("http://api.guildwars2.com/v2/items?ids=")
	for k, v := range ids {
		urlBuilder.WriteString(strconv.Itoa(v))
		if k != len(ids)-1 {
			urlBuilder.WriteString(",")
		}
	}
	url := urlBuilder.String()
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	req.Header.Set("content-type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	responseBody, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		return nil, errReadAll
	}

	if r.StatusCode >= 299 {
		return nil, errors.New(fmt.Sprintf("wrong statusCode: %d, body := ", r.StatusCode))
	}
	baseItemData := make([]gw2models.BaseItemData, len(ids))

	if err := json.Unmarshal(responseBody, &baseItemData); err != nil {
		return nil, err
	}

	return baseItemData, nil
}

func (c Client) GetPrice(id int) (*gw2models.OrderObject, error) {
	r, err := http.Get(fmt.Sprintf("https://api.guildwars2.com/v2/commerce/prices/%d", id))
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.New("Non 200 status code")
	}
	object := gw2models.OrderObject{}
	responseBody, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		return nil, errReadAll
	}

	if err := json.Unmarshal(responseBody, &object); err != nil {
		return nil, err
	}

	return &object, nil
}

func (c Client) GetPriceOfMultiple(ids []int) ([]gw2models.OrderObject, error) {
	if ids == nil || len(ids) == 0 {
		return nil, errors.New("No Ids passed")
	}
	urlBuilder := strings.Builder{}
	urlBuilder.WriteString("https://api.guildwars2.com/v2/commerce/prices?ids=")
	for k, v := range ids {
		urlBuilder.WriteString(strconv.Itoa(v))
		if k != len(ids)-1 {
			urlBuilder.WriteString(",")
		}
	}

	url := urlBuilder.String()

	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}

	r, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	if r.StatusCode > 300 {
		return nil, errors.New("Non 200 status code")
	}
	objs := []gw2models.OrderObject{}
	responseBody, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		return nil, errReadAll
	}

	if err := json.Unmarshal(responseBody, &objs); err != nil {
		return nil, err
	}
	return objs, nil
}

func (c Client) BuyTransactions() ([]gw2models.ActiveTransaction, error) {
	if c.Apikey == "" {
		return nil, errors.New("invalid API KEY")
	}

	url := "https://api.guildwars2.com/v2/commerce/transactions/current/buys"
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Apikey))
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	responseBody, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		return nil, errReadAll
	}

	if r.StatusCode >= 299 {
		return nil, errors.New(fmt.Sprintf("wrong statusCode: %d, body := ", r.StatusCode))
	}
	transactions := []gw2models.ActiveTransaction{}

	if err := json.Unmarshal(responseBody, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil

}

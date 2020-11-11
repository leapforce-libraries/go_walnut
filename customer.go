package walnut

import (
	"fmt"
	"strconv"
	"time"
)

// Customer stores Customer from Walnut
//
type Customer struct {
	UserIdentifier    string   `json:"userIdentifier"`
	PassIdentifier    string   `json:"passIdentifier"`
	PassAdded         string   `json:"passAdded"`
	UserEmail         string   `json:"userEmail"`
	UserName          string   `json:"userName"`
	UserBirthday      string   `json:"userBirthday"`
	UserMobileNumber  string   `json:"userMobileNumber"`
	UserUnsubscribed  bool     `json:"userUnsubscribed"`
	UserUpdated       string   `json:"userUpdated"`
	UserRegistered    string   `json:"userRegistered"`
	UserLicensePlates []string `json:"userLicensePlates"`
}

// GetChanges retrieves changed Customers from Walnut
//
func (w *Walnut) GetChanges(time time.Time) ([]Customer, error) {
	urlStr := "%sstore/%s/changes?date=%s&page=%s"
	page := 0
	rowCount := 1

	err := w.PostLogin()
	if err != nil {
		return nil, err
	}

	customers := []Customer{}

	for rowCount > 0 {
		page++

		layout := "2006-01-02T15:04:05-0700"
		url := fmt.Sprintf(urlStr, w.ApiURL, w.StoreIdentifier, time.Format(layout), strconv.Itoa(page))
		//fmt.Println(url)

		cs := []Customer{}

		err := w.Get(url, &cs)
		if err != nil {
			return nil, err
		}

		for _, c := range cs {
			customers = append(customers, c)
		}

		rowCount = len(cs)
	}

	if len(customers) == 0 {
		customers = nil
	}

	return customers, nil
}

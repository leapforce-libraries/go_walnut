package walnut

import (
	"fmt"
	"strconv"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Customer stores Customer from Service
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

// GetChanges retrieves changed Customers from Service
//
func (service *Service) GetChanges(time time.Time) ([]Customer, *errortools.Error) {
	urlStr := "%s/store/%s/changes?date=%s&page=%s"
	page := 0
	rowCount := 1

	if !service.static {
		e := service.PostLogin()
		if e != nil {
			return nil, e
		}
	}

	customers := []Customer{}

	for rowCount > 0 {
		page++

		layout := "2006-01-02T15:04:05-0700"
		url := fmt.Sprintf(urlStr, APIURL, service.StoreIdentifier, time.Format(layout), strconv.Itoa(page))
		//fmt.Println(url)

		cs := []Customer{}

		e := service.Get(url, &cs)
		if e != nil {
			return nil, e
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

package walnut

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
)

// Login stores Login from Service
//
type Login struct {
	AccountToken    string `json:"accountToken"`
	StoreIdentifier string `json:"storeIdentifier"`
}

// GetChanges retrieves changed Customers from Service
//
func (service *Service) PostLogin() *errortools.Error {
	urlStr := "%s/login"
	url := fmt.Sprintf(urlStr, APIURL)
	//fmt.Println(url)

	login := Login{}

	data := make(map[string]string)
	data["accountEmailAddress"] = service.EmailAddress
	data["accountPassword"] = service.Password
	data["partnerToken"] = service.PartnerToken

	e := service.Post(url, data, &login, false, false)
	if e != nil {
		return e
	}

	service.StoreIdentifier = login.StoreIdentifier
	service.AccountToken = login.AccountToken

	//fmt.Println(service)

	return nil
}

package walnut

import (
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
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
	login := Login{}

	bodyModel := struct {
		AccountEmailAddress string `json:"accountEmailAddress"`
		AccountPassword     string `json:"accountPassword"`
		PartnerToken        string `json:"partnerToken"`
	}{
		service.emailAddress,
		service.password,
		service.partnerToken,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		URL:           service.url("login"),
		BodyModel:     bodyModel,
		ResponseModel: &login,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	service.storeIdentifier = login.StoreIdentifier
	service.accountToken = login.AccountToken

	return nil
}

package walnut

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

const (
	APIURL string = "https://walnutbackend.com/api/v1"
)

// type
//
type Service struct {
	EmailAddress    string
	Password        string
	PartnerToken    string
	StoreIdentifier string
	AccountToken    string
	static          bool
	isLive          bool
}

// Response represents highest level of exactonline api response
//
type Response struct {
	Results *json.RawMessage `json:"results"`
}

func NewService(emailAddress string, password string, partnerToken string, isLive bool) (*Service, *errortools.Error) {
	service := new(Service)

	if emailAddress == "" {
		return nil, errortools.ErrorMessage("Service emailAddress not provided")
	}
	if password == "" {
		return nil, errortools.ErrorMessage("Service password not provided")
	}
	if partnerToken == "" {
		return nil, errortools.ErrorMessage("Service partnerToken not provided")
	}

	service.EmailAddress = emailAddress
	service.Password = password
	service.PartnerToken = partnerToken
	service.isLive = isLive
	service.static = false

	return service, nil
}

func NewServiceStatic(storeIdenitifier string, accountToken string, isLive bool) (*Service, *errortools.Error) {
	service := new(Service)

	if storeIdenitifier == "" {
		return nil, errortools.ErrorMessage("Service StoreIdenitifier not provided")
	}
	if accountToken == "" {
		return nil, errortools.ErrorMessage("Service AccountToken not provided")
	}

	service.StoreIdentifier = storeIdenitifier
	service.AccountToken = accountToken
	service.isLive = isLive
	service.static = true

	return service, nil
}

// Get is a generic Get method
//
func (service *Service) Get(url string, model interface{}) *errortools.Error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errortools.ErrorMessage(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("WalnutPass %s", service.AccountToken))

	attempts := 10
	attempt := 1

	res := new(http.Response)

	for attempt < attempts {
		// Send out the HTTP request
		res, err = client.Do(req)
		if err != nil {
			attempt++
			fmt.Println("url:", url)
			fmt.Println("error:", err.Error())
			fmt.Println("starting attempt:", attempt)

			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	response := Response{}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	err = json.Unmarshal(*response.Results, &model)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	return nil
}

// Post is a generic Post method
//
func (service *Service) Post(url string, values map[string]string, model interface{}, authorize bool, responseWrapped bool) *errortools.Error {
	client := &http.Client{}

	buf := new(bytes.Buffer)
	if values != nil {
		json.NewEncoder(buf).Encode(values)
	} else {
		buf = nil
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	// add headers
	req.Header.Set("accept", "application/json")
	if authorize {
		req.Header.Set("authorization", fmt.Sprintf("WalnutPass %s", service.AccountToken))
	}

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	if responseWrapped {

		response := Response{}

		err = json.Unmarshal(b, &response)
		if err != nil {
			return errortools.ErrorMessage(err)
		}

		b = *response.Results
	}

	err = json.Unmarshal(b, &model)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	return nil
}

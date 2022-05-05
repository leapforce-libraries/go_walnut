package walnut

import (
	"encoding/json"
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName    string = "Walnut"
	apiUrl     string = "https://walnutbackend.com/api/v1"
	dateLayout string = "2006-01-02T15:04:05-0700"
)

type Response struct {
	Results *json.RawMessage `json:"results"`
}

type Service struct {
	emailAddress    string
	password        string
	partnerToken    string
	storeIdentifier string
	accountToken    string
	static          bool
	httpService     *go_http.Service
}

type ServiceConfig struct {
	EmailAddress string
	Password     string
	PartnerToken string
}

func NewService(config *ServiceConfig) (*Service, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if config.EmailAddress == "" {
		return nil, errortools.ErrorMessage("Service emailAddress not provided")
	}
	if config.Password == "" {
		return nil, errortools.ErrorMessage("Service password not provided")
	}
	if config.PartnerToken == "" {
		return nil, errortools.ErrorMessage("Service partnerToken not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		emailAddress: config.EmailAddress,
		password:     config.Password,
		partnerToken: config.PartnerToken,
		static:       false,
		httpService:  httpService,
	}, nil
}

type ServiceStaticConfig struct {
	StoreIdentifier string
	AccountToken    string
}

func NewServiceStatic(config *ServiceStaticConfig) (*Service, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if config.StoreIdentifier == "" {
		return nil, errortools.ErrorMessage("Service StoreIdentifier not provided")
	}
	if config.AccountToken == "" {
		return nil, errortools.ErrorMessage("Service AccountToken not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		storeIdentifier: config.StoreIdentifier,
		accountToken:    config.AccountToken,
		static:          true,
		httpService:     httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add token
	header := http.Header{}
	header.Set("authorization", fmt.Sprintf("WalnutPass %s", service.accountToken))
	(*requestConfig).NonDefaultHeaders = &header

	request, response, e := service.httpService.HttpRequest(requestConfig)

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiUrl, path)
}

func (service *Service) getResponse(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	responseModel := requestConfig.ResponseModel

	response := Response{}
	requestConfig.ResponseModel = &response

	req, res, e := service.httpRequest(requestConfig)
	if e != nil {
		return req, res, e
	}

	err := json.Unmarshal(*response.Results, responseModel)
	if err != nil {
		return req, res, errortools.ErrorMessage(err)
	}

	return req, res, nil
}

func (service Service) APIName() string {
	return apiName
}

func (service Service) APIKey() string {
	return service.accountToken
}

func (service Service) APICallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) APIReset() {
	service.httpService.ResetRequestCount()
}

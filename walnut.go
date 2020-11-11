package walnut

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	types "github.com/leapforce-libraries/go_types"
)

// type
//
type Walnut struct {
	ApiURL          string
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

func New(apiURL string, emailAddress string, password string, partnerToken string, isLive bool) (*Walnut, error) {
	w := new(Walnut)

	if apiURL == "" {
		return nil, &types.ErrorString{"Walnut ApiUrl not provided"}
	}
	if emailAddress == "" {
		return nil, &types.ErrorString{"Walnut emailAddress not provided"}
	}
	if password == "" {
		return nil, &types.ErrorString{"Walnut password not provided"}
	}
	if partnerToken == "" {
		return nil, &types.ErrorString{"Walnut partnerToken not provided"}
	}

	w.ApiURL = apiURL
	w.EmailAddress = emailAddress
	w.Password = password
	w.PartnerToken = partnerToken
	w.isLive = isLive
	w.static = false

	if !strings.HasSuffix(w.ApiURL, "/") {
		w.ApiURL = w.ApiURL + "/"
	}

	return w, nil
}

func NewStatic(apiURL string, storeIdenitifier string, accountToken string, isLive bool) (*Walnut, error) {
	w := new(Walnut)

	if apiURL == "" {
		return nil, &types.ErrorString{"Walnut ApiUrl not provided"}
	}
	if storeIdenitifier == "" {
		return nil, &types.ErrorString{"Walnut StoreIdenitifier not provided"}
	}
	if accountToken == "" {
		return nil, &types.ErrorString{"Walnut AccountToken not provided"}
	}

	w.ApiURL = apiURL
	w.StoreIdentifier = storeIdenitifier
	w.AccountToken = accountToken
	w.isLive = isLive
	w.static = true

	if !strings.HasSuffix(w.ApiURL, "/") {
		w.ApiURL = w.ApiURL + "/"
	}

	return w, nil
}

// Get is a generic Get method
//
func (w *Walnut) Get(url string, model interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("WalnutPass %s", w.AccountToken))

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
		return err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	response := Response{}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*response.Results, &model)
	if err != nil {
		return err
	}

	return nil
}

// Post is a generic Post method
//
func (w *Walnut) Post(url string, values map[string]string, model interface{}, authorize bool, responseWrapped bool) error {
	client := &http.Client{}

	buf := new(bytes.Buffer)
	if values != nil {
		json.NewEncoder(buf).Encode(values)
	} else {
		buf = nil
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}

	// add headers
	req.Header.Set("accept", "application/json")
	if authorize {
		req.Header.Set("authorization", fmt.Sprintf("WalnutPass %s", w.AccountToken))
	}

	// Send out the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	if responseWrapped {

		response := Response{}

		err = json.Unmarshal(b, &response)
		if err != nil {
			return err
		}

		b = *response.Results
	}

	err = json.Unmarshal(b, &model)
	if err != nil {
		return err
	}

	return nil
}

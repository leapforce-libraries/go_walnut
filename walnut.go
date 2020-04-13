package walnut

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	types "github.com/Leapforce-nl/go_types"
)

// type
//
type Walnut struct {
	ApiURL           string
	StoreIdenitifier string
	AccountToken     string
	IsLive           bool
}

// Response represents highest level of exactonline api response
//
type Response struct {
	Results *json.RawMessage `json:"results"`
}

func New(apiURL string, storeIdenitifier string, accountToken string, isLive bool) (*Walnut, error) {
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
	w.StoreIdenitifier = storeIdenitifier
	w.AccountToken = accountToken
	w.IsLive = isLive

	if !strings.HasSuffix(w.ApiURL, "/") {
		w.ApiURL = w.ApiURL + "/"
	}

	return w, nil
}

// generic Get method
//
func (w *Walnut) Get(url string, model interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("WalnutPass %s", w.AccountToken))

	// Send out the HTTP request
	res, err := client.Do(req)
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

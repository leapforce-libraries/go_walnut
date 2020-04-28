package walnut

import (
	"fmt"
)

// Login stores Login from Walnut
//
type Login struct {
	AccountToken    string `json:"accountToken"`
	StoreIdentifier string `json:"storeIdentifier"`
}

// GetChanges retrieves changed Customers from Walnut
//
func (w *Walnut) PostLogin() error {
	urlStr := "%slogin"
	url := fmt.Sprintf(urlStr, w.ApiURL)
	//fmt.Println(url)

	login := Login{}

	data := make(map[string]string)
	data["accountEmailAddress"] = w.EmailAddress
	data["accountPassword"] = w.Password
	data["partnerToken"] = w.PartnerToken

	err := w.Post(url, data, &login, false, false)
	if err != nil {
		return err
	}

	w.StoreIdentifier = login.StoreIdentifier
	w.AccountToken = login.AccountToken

	//fmt.Println(w)

	return nil
}

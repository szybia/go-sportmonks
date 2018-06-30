package sportmonks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var apiToken = ""
var apiURL = "https://soccer.sportmonks.com/api/v2.0/"

//SetAPIToken quite self-explanatory.
//APIParameters specifies the options supplied to the Get function
type APIParameters struct {
	Endpoint string
	Include  string
	Page     int
	AllPages bool
}

//NewAPIParameters uses
func NewAPIParameters(endpoint string, include string, page int, allPages bool) *APIParameters {
	A := APIParameters{
		Endpoint: endpoint,
		Include:  include,
		Page:     NoPageDefault,
		AllPages: allPages}

	if page != NoPageDefault {
		A.Page = page
	}
	return &A
}
func SetAPIToken(s string) {
	if len(s) != 0 {
		apiToken = s
	}
}

//Get API request.
func Get(endpoint string) map[string]string {
	payload := map[string]string{"api_token": apiToken}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Print(err)
	}

	q := req.URL.Query()
	q.Add("api_token", apiToken)
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		log.Print(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	var objmap map[string]interface{}
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		log.Print(err)
	}

	fmt.Println(req.URL.String())
	return payload
}

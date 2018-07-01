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

//NoPageDefault specifies the default when a specific page is not requested
var NoPageDefault = -1

//SetAPIToken sets the API token for sportmonks
func SetAPIToken(s string) {
	if len(s) != 0 {
		apiToken = s
	}
}

//Get API request
func Get(endpoint string, include string, page int, allPages bool) ([]byte, error) {
	if endpoint == "" {
		return []byte{}, errors.New("no endpoint provided")
	} else if apiToken == "" {
		return []byte{}, errors.New("apiToken has not been set")
	}

	requestURL := apiURL + endpoint
	r, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return []byte{}, err
	}

	q := r.URL.Query()
	q.Add("api_token", apiToken)
	if include != "" {
		q.Add("include", include)
	}
	if page != FirstPage {
		q.Add("page", strconv.Itoa(page))
		allPages = false
	}
	r.URL.RawQuery = q.Encode()

	resp, err := http.Get(r.URL.String())
	if err != nil {
		return []byte{}, err
	}

	body, err := jason.NewObjectFromReader(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	data, err := body.GetObjectArray("data")
	if err != nil {
		return []byte{}, err
	}

	if allPages {
		pages, err := body.GetInt64("meta", "pagination", "total_pages")
		//	No error means endpoint is paginated
		if err == nil {
			if pages > 1 {
				c := make(chan paginatedRequest)
				requests := make([][]*jason.Object, pages)
				for i := int64(2); i <= pages; i++ {
					go getRequest(r.URL.String(), i, c)
				}

				for i := int64(2); i <= pages; i++ {
					g := <-c
					requests[g.pageNumber-1] = g.data
				}

				for i := int64(1); i < pages; i++ {
					data = append(data, requests[i]...)
				}

			}
		}
	}
	m, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}

	return m, nil
}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Fatal(err)
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

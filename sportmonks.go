package sportmonks

import (
	"net/http"
	"os"
)

var apiURL = "https://soccer.sportmonks.com/api/v2.0/"

//Test http request
func Test() {
	response, err := http.Get(apiURL)
	var responseByte []byte
	if err != nil {
		os.Exit(1)
	}
	_, err = response.Body.Read(responseByte)
	if err != nil {
		os.Exit(0)
	}
}

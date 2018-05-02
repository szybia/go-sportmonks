package sportmonks

import (
	"net/http"
	"os"
)

var apiToken = ""
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

func setAPIToken(s string) {
	if len(s) != 0 {
		apiToken = s
	}
}

func get(endpoint string) map[string]string {
	payload := map[string]string{"api_token": apiToken}
	return payload
}

# golang-sportmonks
Golang wrapper for the Sportmonks Soccer API  
#### Pull requests and recommendations are welcome!

## Installation
```bash
$ go get github.com/BialkowskiSz/go-sportmonks
```

## Usage
```golang
import (
	"fmt"
	"log"

	sportmonks "github.com/BialkowskiSz/go-sportmonks"
)

func main() {
	sportmonks.SetAPIToken("<YOUR_TOKEN_HERE>")
	g, err := sportmonks.Get("fixtures/between/2016-01-01/2018-01-01", "", 0, false)
  //  Can also use sportmonks globals for code clarity and readability
  //  g, err := sportmonks.Get("fixtures/between/2016-01-01/2018-01-01", sportmonks.NoIncludes, sportmonks.FirstOrAllPages, sportmonks.SinglePage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(g))
}
```

- Contains functions for all Sportmonks endpoints as well as base custom Get function
```golang
g, err := sportmonks.Leagues("", 0, true)
g, err := sportmonks.Seasons("", 0, true)
g, err := sportmonks.LivescoresNow(sportmonks.NoIncludes)
```
- Contains build in logger which is used for goroutine errors. Logger can be altered to requirements
```golang
	f, err := os.OpenFile("sportmonks_log.txt", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	sportmonks.Logger.SetOutput(f)
	//Code...
```

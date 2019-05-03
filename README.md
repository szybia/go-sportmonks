[![Build Status](https://travis-ci.org/szybia/go-sportmonks.svg?branch=master)](https://travis-ci.org/szybia/go-sportmonks)
[![Go Report Card](https://goreportcard.com/badge/github.com/szybia/go-sportmonks)](https://goreportcard.com/report/github.com/szybia/go-sportmonks)
[![License](https://img.shields.io/badge/license-apache2.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/szybia/go-sportmonks/master/LICENSE)

# golang-sportmonks
Golang wrapper for the Sportmonks Soccer API  
#### Pull requests are welcome!

## Installation
```bash
$ go get github.com/szybia/go-sportmonks
```

## Usage
```golang
package main

import (
	"fmt"
	"log"

	sm "github.com/szybia/go-sportmonks"
)

func main() {
	sm.SetAPIToken("<YOUR_TOKEN_HERE>")
	g, err := sm.Get("fixtures/between/2016-01-01/2018-01-01", "", 0, false)
	//  Can also use sportmonks globals for code clarity and readability
	//  g, err := sm.Get("fixtures/between/2016-01-01/2018-01-01", sm.NoIncludes, sm.FirstOrAllPages, sm.SinglePage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(g))
}
```

- Contains functions for all Sportmonks endpoints as well as base custom Get function
```golang
g, err := sm.Leagues("", 0, true)
g, err := sm.Seasons("", 0, true)
g, err := sm.LivescoresNow(sm.NoIncludes)
```
- Contains build in logger which is used for goroutine errors. Logger can be altered to requirements
```golang
f, err := os.OpenFile("sportmonks_log.txt", os.O_CREATE|os.O_APPEND, 0666)
if err != nil {
	log.Fatal(err)
}
sm.Logger.SetOutput(f)
//Code...

```

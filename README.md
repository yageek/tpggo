# tpggo

An http client for the TPG Open Data.

For more explanation: [TPG Open Data](http://www.tpg.ch/fr/web/open-data/mode-d-emploi)

## Installation

```shell
go get github.com/gophersch/tpggo
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/gophersch/tpggo"
)

const (
	apiKey = "MY_KEY"
)

func main() {

	client := tpggo.NewClient(apiKey)

	// Get all the stops
	if resp, err := client.GetStops(); err != nil {
		panic(err)
	} else {
		for _, stop := range resp.Stops {
			fmt.Printf("Stop Name: %s | Code: %s \n", stop.Code, stop.Name)
		}
	}
}
```

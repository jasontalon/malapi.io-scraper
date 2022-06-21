package main

import (
	"fmt"
	"github.com/jasontalon/malapi.io-scraper"
)

func main() {

	items, err := malapi.Get()

	if err != nil {
		panic(err)
	}

	err = malapi.ExportToCsv(&items)

	if err != nil {
		panic(err)
	}

	fmt.Println("done")
	fmt.Println("results exported to malapi.csv file")
}

package malapi

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

const (
	url = "https://malapi.io"
)

type Api struct {
	Name              string `json:"Function Name"`
	Description       string
	Library           string
	AssociatedAttacks []string `json:"Associated Attacks"`
	Documentation     string
	Created           string
	LastUpdate        string
	Credits           string
}

func Get() {
	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var mals []string

	f := func(i int, s *goquery.Selection) {
		mals = append(mals, s.Text())
	}

	doc.Find("a[href*='/winapi']").Each(f)

	var apis []Api
	for _, mal := range mals {
		api, err := GetByName(mal)
		if err == nil {
			apis = append(apis, api)
		}
	}

	fmt.Println("done!")
}

func GetByName(name string) (api Api, err error) {
	response, err := http.Get(url + "/winapi/" + name)

	if err != nil {
		return api, err
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return api, err
	}
	dict := make(map[string]interface{})

	doc.Find(".heading").Each(func(i int, selection *goquery.Selection) {
		header := strings.TrimSpace(selection.Text())
		value := strings.TrimSpace(selection.Next().Text())
		if header == "Associated Attacks" {
			dict[header] = strings.Fields(value)
		} else {
			dict[header] = value
		}
	})

	doc.Find(".detail-container > .square-box > .content").Each(func(i int, selection *goquery.Selection) {
		s := strings.Split(selection.Text(), ":")
		header := strings.ReplaceAll(s[0], " ", "")
		value := s[len(s)-1]
		dict[header] = value
	})
	fmt.Println(dict)

	data, err := json.Marshal(dict)

	if err != nil {
		return api, err
	}

	json.Unmarshal(data, &api)

	return api, nil
}

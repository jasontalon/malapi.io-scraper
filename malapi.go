package malapi

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	url         = "https://malapi.io"
	workerCount = 4
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

func ExportToCsv(list *[]Api) (err error) {
	f, err := os.Create("malapi.csv")
	if err != nil {
		return err
	}
	defer func(f *os.File) error {
		err := f.Close()
		if err != nil {
			return err
		}
		return nil
	}(f)

	writer := csv.NewWriter(f)

	header := []string{
		"Name",
		"Description",
		"Library",
		"AssociatedAttacks",
		"Documentation",
		"Created",
		"LastUpdate",
		"Credits",
	}
	var records [][]string
	records = append(records, header)
	for _, item := range *list {
		records = append(records, []string{
			item.Name,
			item.Description,
			item.Library,
			strings.Join(item.AssociatedAttacks, ","),
			item.Documentation,
			item.Created,
			item.LastUpdate,
			item.Credits,
		})
	}

	err = writer.WriteAll(records)

	if err != nil {
		return err
	}

	return nil
}

func Get() (list []Api, err error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var names []string

	doc.Find("a[href*='/winapi']").Each(func(i int, s *goquery.Selection) {
		names = append(names, s.Text())
	})

	chunkedNames := chunk(names, 10)

	workChunk := chunk(func() []int {
		s := make([]int, len(chunkedNames))
		for i := range s {
			s[i] = i
		}
		return s
	}(), workerCount)

	for _, o := range workChunk {
		ch := make(chan []Api, len(o))
		for _, i := range o {
			go readDetails(chunkedNames[i], ch)
		}

		for i := 0; i < len(o); i++ {
			select {
			case d := <-ch:
				list = append(list, d...)
			}
		}
	}

	return list, nil
}

func readDetails(names []string, ch chan []Api) {
	var l []Api
	for _, name := range names {
		api, err := readDetail(name)
		if err == nil {
			l = append(l, api)
		}
	}
	ch <- l
}

func readDetail(name string) (api Api, err error) {
	fmt.Println("reading " + name)
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

	data, err := json.Marshal(dict)

	if err != nil {
		return api, err
	}

	err = json.Unmarshal(data, &api)

	if err != nil {
		return api, err
	}

	return api, nil
}

// Chunk https://github.com/samber/lo/blob/305f565368f9959501e17628d644101fd18d4de1/slice.go#L139
// Chunk returns an array of elements split into groups the length of size. If array can't be split evenly,
// the final chunk will be the remaining elements.
func chunk[T any](collection []T, size int) [][]T {
	if size <= 0 {
		panic("Second parameter must be greater than 0")
	}

	result := make([][]T, 0, len(collection)/2+1)
	length := len(collection)

	for i := 0; i < length; i++ {
		chunk := i / size

		if i%size == 0 {
			result = append(result, make([]T, 0, size))
		}

		result[chunk] = append(result[chunk], collection[i])
	}

	return result
}

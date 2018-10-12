package surfline

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func requestQuery(keyword string) (*queryResp, error) {

	baseUrl := "https://services.surfline.com/search/site"

	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		logrus.Errorf("Can't create HTTP Request to %s Error", baseUrl)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("q", keyword)
	req.URL.RawQuery = q.Encode()

	var client = &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		logrus.Errorf("Error Occured send HTTP Request to", req.URL.String())
		return nil, err
	}
	defer res.Body.Close()

	//parse response JSON data
	respJson := new(queryResp)
	if err := json.NewDecoder(res.Body).Decode(respJson); err != nil {
		logrus.Error("Faile to decode query response")
		return nil, err
	}
	return respJson, nil
}

type Item struct {
	Title    string
	Url      string
	Type     string
	SubTitle string
}

func newItem(s *sourceStruct) *Item {
	it := new(Item)
	it.Title = s.Name
	it.Url = s.Href
	it.SubTitle = strings.Join(s.BreadCrumbs, "/")
	return it
}

func Query(keyword string) ([]*Item, error) {
	resp, err := requestQuery(keyword)
	if err != nil {
		return nil, err
	}

	//to remove duplicated spot info
	urlSet := map[string]bool{}
	var items []*Item

	checkAndAdd := func(s *sourceStruct) {
		href := s.Href
		if _, exist := urlSet[href]; exist {
			return
		}
		urlSet[href] = true
		it := newItem(s)
		items = append(items, it)
	}

	for _, it := range *resp {
		for _, spotSuggest := range it.Suggest.SpotSuggest {
			for _, options := range spotSuggest.Options {
				checkAndAdd(&options.Source)
			}
		}

		for _, hits := range it.Hits.Hits {
			checkAndAdd(&hits.Source)
		}
	}
	return items, nil
}

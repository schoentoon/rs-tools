package info

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TravellingMerchantShop struct {
	Products []TravellingMerchantProduct
}

type TravellingMerchantProduct struct {
	Name        string
	Cost        int
	Quantity    string
	Description string
}

func parseProduct(s *goquery.Selection) (out TravellingMerchantProduct, err error) {
	s.Find("td").Each(func(i int, ss *goquery.Selection) {
		switch i {
		case 1:
			out.Name = ss.Text()
		case 2:
			out.Cost, err = strconv.Atoi(strings.ReplaceAll(ss.Text(), ",", ""))
		case 3:
			out.Quantity = ss.Text()
		case 4:
			out.Description = ss.Text()
		}
	})
	return
}

func TravellingMerchant(client *http.Client) (*TravellingMerchantShop, error) {
	params := url.Values{
		"action":             {"parse"},
		"format":             {"json"},
		"text":               {"{{Travelling Merchant}}"},
		"contentmodel":       {"wikitext"},
		"prop":               {"text"},
		"disablelimitreport": {"1"},
	}
	resp, err := client.Get(fmt.Sprintf("https://runescape.wiki/api.php?%s", params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Status: %d %s", resp.StatusCode, resp.Status)
	}

	wrapper := struct {
		Parse struct {
			Text struct {
				Data string `json:"*"`
			} `json:"text"`
		} `json:"parse"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(wrapper.Parse.Text.Data))
	if err != nil {
		return nil, err
	}

	out := &TravellingMerchantShop{}

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 || err != nil {
			return
		}
		product, e := parseProduct(s)
		if e != nil {
			err = e
			return
		}
		out.Products = append(out.Products, product)
	})

	return out, err
}

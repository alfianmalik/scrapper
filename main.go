package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Shops struct {
	ShopUrl  string `json:"shopUrl"`
	Products Products
}

type Products struct {
	ProductName string `json:"productName"`
	ProductUrl  string `json:"productUrl"`
	Price       int    `json:"price"`
}

func main() {
	escaped := regexp.QuoteMeta("tokopedia.com")
	r := regexp.MustCompile(`^https?:\/\/[a-z]*\.?` + escaped + `.*`)
	c := colly.NewCollector(
		colly.URLFilters(r),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"),
	)
	c.SetRequestTimeout(120 * time.Second)

	shops := make([]Shops, 0)
	c.OnHTML("div.pcv3__container", func(e *colly.HTMLElement) {
		e.ForEach("div.css-974ipl", func(i int, h *colly.HTMLElement) {
			reg, err := regexp.Compile("[^0-9]+")
			if err != nil {
				log.Fatal(err)
			}
			processedString := reg.ReplaceAllString(e.ChildText("div.css-1ksb19c"), "")
			price, err := strconv.Atoi(processedString)
			if err != nil {
				log.Fatal(err)
			}

			shopName := strings.Split(e.ChildAttr("a.css-gwkf0u", "href"), "/")

			shop := Shops{
				ShopUrl: "https://www.tokopedia.com/" + strings.ToLower(strings.ReplaceAll(shopName[3], " ", "")),
				Products: Products{
					ProductName: e.ChildText("div.css-1b6t4dn"),
					ProductUrl:  e.ChildAttr("a.css-gwkf0u", "href"),
					Price:       price,
				},
			}
			shops = append(shops, shop)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		file, err := json.MarshalIndent(shops, "", " ")
		if err != nil {
			log.Println("Unable to create json file")
			return
		}
		// fmt.Println("Writing data to file")
		if err = ioutil.WriteFile("shops.json", file, 0644); err == nil {
			// fmt.Println("Data written to file successfully")
		}
	})

	products := []string{"Zjiang ZJ-5890T", "Sennheiser HD 202", "sepatu balerina hitam flat"}
	for i := 0; i < len(products); i++ {
		link := fmt.Sprintf("https://www.tokopedia.com/search?st=product&q=%s", url.QueryEscape(products[i]))
		c.Visit(link)
	}

}

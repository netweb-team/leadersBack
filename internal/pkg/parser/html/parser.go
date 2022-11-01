package html

import (
	"encoding/json"
	"fmt"
	"leaders_apartments/internal/pkg/domain"
	"net/http"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/html"
)

const (
	params         = "&p=%d"
	psize          = 28
	coordinatesReq = "https://www.cian.ru/ajax/map/roundabout/?engine_version=2&deal_type=sale&offer_type=flat&region=1&minprice=%d&maxprice=%d&room2=1"
)

func Search(url string) *domain.SearchPage {
	links := make([]string, 0)
	lastPage := 1
	doc := ListParse(&links, url, lastPage)
	if doc == nil {
		return nil
	}

	//expr, _ := xpath.Compile(`count(//a[@class="pagination-block__page-link"])`)
	//v := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(float64)
	ads, _ := htmlquery.Query(doc, `//div[@data-name="SummaryHeader"]`)
	adsTxt := htmlquery.InnerText(ads)
	lastPage, _ = strconv.Atoi(adsTxt[strings.Index(adsTxt, " ")+1 : strings.LastIndex(adsTxt, " ")])
	lastPage /= psize

	for i := 2; i <= lastPage; i++ {
		if err := ListParse(&links, url, i); err == nil {
			return nil
		}
	}
	return &domain.SearchPage{Links: links, Count: len(links)}
}

func ListParse(links *[]string, url string, page int) *html.Node {
	doc, err := htmlquery.LoadURL(url + fmt.Sprintf(params, page))
	if err != nil {
		log.Error("Bad url:", err)
		return nil
	}
	if list, err := htmlquery.QueryAll(doc, `//div[@data-name="LinkArea"]/a`); err == nil {
		for _, el := range list {
			*links = append(*links, htmlquery.SelectAttr(el, "href"))
		}
	} else {
		log.Error(err)
		return nil
	}
	return doc
}

func Ad(url string) *domain.AdPage {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Error("Bad url:", err)
		return nil
	}
	priceNode, err := htmlquery.Query(doc, `//span[@itemprop="price"]`)
	if err != nil {
		log.Error(err)
		return nil
	}
	priceStr := htmlquery.SelectAttr(priceNode, "content")
	price, _ := strconv.Atoi(strings.ReplaceAll(priceStr[:strings.LastIndex(priceStr, " ")], " ", ""))

	resp, err := http.Get(fmt.Sprintf(coordinatesReq, price, price))
	if err != nil {
		log.Error(err)
		return nil
	}
	defer resp.Body.Close()
	coords := new(domain.Coordinates)
	json.NewDecoder(resp.Body).Decode(coords)

	geoNode, _ := htmlquery.Query(doc, `//div[@data-name="Geo"]`)
	addrNode, err := htmlquery.Query(geoNode, `//span[@itemprop="name"]`)
	if err != nil {
		log.Error(err)
		return nil
	}
	addr := htmlquery.SelectAttr(addrNode, "content")
	metroNode, err := htmlquery.Query(geoNode, `//li[position()=1]`)
	if err != nil {
		log.Error(err)
		return nil
	}

	result := &domain.AdPage{
		Address: addr,
		Price:   priceStr,
		Metro:   htmlquery.InnerText(metroNode),
	}
	for k, v := range coords.Data.Points {
		if strings.Contains(addr, v.Content.Address) {
			i := strings.Index(k, " ")
			result.Latitude, result.Longitude = k[:i], k[i+1:]
		}
	}
	summary, _ := htmlquery.QueryAll(doc, `//div[@data-testid="object-summary-description-info"]`)
	for _, el := range summary {
		title, _ := htmlquery.Query(el, `//div[@data-testid="object-summary-description-title"]`)
		value, _ := htmlquery.Query(el, `//div[@data-testid="object-summary-description-value"]`)
		switch htmlquery.InnerText(title) {
		case "Общая":
			result.TotalArea = htmlquery.InnerText(value)
		case "Кухня":
			result.KitchenArea = htmlquery.InnerText(value)
		case "Этаж":
			floor := htmlquery.InnerText(value)
			result.Floor = floor[:strings.Index(floor, " ")]
		case "Построен":
			result.Year = htmlquery.InnerText(value)
		}
	}
	addFeature, _ := htmlquery.QueryAll(doc, `//li[@data-name="AdditionalFeatureItem"]`)
	for _, el := range addFeature {
		title, _ := htmlquery.Query(el, `//span[position()=1]`)
		value, _ := htmlquery.Query(el, `//span[position()=2]`)
		switch htmlquery.InnerText(title) {
		case "Балкон/лоджия":
			result.Balcony = htmlquery.InnerText(value)
		case "Ремонт":
			result.Renovation = htmlquery.InnerText(value)
		}
	}
	items, _ := htmlquery.QueryAll(doc, `//div[@data-name="Item"]`)
	for _, el := range items {
		title, _ := htmlquery.Query(el, `//div[position()=1]`)
		value, _ := htmlquery.Query(el, `//div[position()=2]`)
		switch htmlquery.InnerText(title) {
		case "Год постройки":
			result.Year = htmlquery.InnerText(value)
			// case "Тип дома":
			// 	result.Year = htmlquery.InnerText(value)
		}
	}
	return result
}

package html

import (
	"encoding/json"
	"fmt"
	"leaders_apartments/internal/pkg/domain"
	"net/http"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
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

	ads := htmlquery.QuerySelector(doc, divSummary)
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
	list := htmlquery.QuerySelectorAll(doc, divLinkArea)
	for _, el := range list {
		*links = append(*links, htmlquery.SelectAttr(el, "href"))
	}
	return doc
}

func Ad(url string) *domain.AdPage {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Error("Bad url:", err)
		return nil
	}
	priceNode := htmlquery.QuerySelector(doc, spanPrice)
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

	geoNode := htmlquery.QuerySelector(doc, divGeo)
	addrNode := htmlquery.QuerySelector(geoNode, spanName)
	addr := htmlquery.SelectAttr(addrNode, "content")
	metroNode := htmlquery.QuerySelector(geoNode, li1)

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
	summary := htmlquery.QuerySelectorAll(doc, divInfo)
	for _, el := range summary {
		title := htmlquery.QuerySelector(el, divTitle)
		value := htmlquery.QuerySelector(el, divValue)
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
	addFeature := htmlquery.QuerySelectorAll(doc, liFeature)
	for _, el := range addFeature {
		title := htmlquery.QuerySelector(el, span1)
		value := htmlquery.QuerySelector(el, span2)
		switch htmlquery.InnerText(title) {
		case "Балкон/лоджия":
			result.Balcony = htmlquery.InnerText(value)
		case "Ремонт":
			result.Renovation = htmlquery.InnerText(value)
		}
	}
	items := htmlquery.QuerySelectorAll(doc, divItem)
	for _, el := range items {
		title := htmlquery.QuerySelector(el, div1)
		value := htmlquery.QuerySelector(el, div2)
		switch htmlquery.InnerText(title) {
		case "Год постройки":
			result.Year = htmlquery.InnerText(value)
			// case "Тип дома":
			// 	result.Year = htmlquery.InnerText(value)
		}
	}
	return result
}

var divSummary, divLinkArea, spanPrice, divGeo, spanName, li1, divInfo, divTitle, divValue, liFeature, span1, span2, div1, div2, divItem *xpath.Expr

func init() {
	divSummary = xpath.MustCompile(`//div[@data-name="SummaryHeader"]`)
	divLinkArea = xpath.MustCompile(`//div[@data-name="LinkArea"]/a`)
	spanPrice = xpath.MustCompile(`//span[@itemprop="price"]`)
	divGeo = xpath.MustCompile(`//div[@data-name="Geo"]`)
	spanName = xpath.MustCompile(`//span[@itemprop="name"]`)
	li1 = xpath.MustCompile(`//li[position()=1]`)
	divInfo = xpath.MustCompile(`//div[@data-testid="object-summary-description-info"]`)
	divTitle = xpath.MustCompile(`//div[@data-testid="object-summary-description-title"]`)
	divValue = xpath.MustCompile(`//div[@data-testid="object-summary-description-value"]`)
	liFeature = xpath.MustCompile(`//li[@data-name="AdditionalFeatureItem"]`)
	span1 = xpath.MustCompile(`//span[position()=1]`)
	span2 = xpath.MustCompile(`//span[position()=2]`)
	div1 = xpath.MustCompile(`//div[position()=1]`)
	div2 = xpath.MustCompile(`//div[position()=2]`)
	divItem = xpath.MustCompile(`//div[@data-name="Item"]`)
}

package html

import (
	"encoding/json"
	"fmt"
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/domain"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	geo "github.com/kellydunn/golang-geo"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/html"
)

const (
	params    = "&p=%d"
	psize     = 28
	yearOld   = 1983
	bywalk    = " пешком"
	minutelen = 16
	unispace  = "\u00a0"
)

func Search(lat, lng float64, segment, url string, room int) []*domain.Row {
	lastPage := 1
	page, doc := listParse(url, lastPage)
	if doc == nil {
		return nil
	}

	adscnt := htmlquery.QuerySelector(doc, divSummary)
	if adscnt != nil {
		adsTxt := htmlquery.InnerText(adscnt)
		lastPage, _ = strconv.Atoi(adsTxt[strings.Index(adsTxt, " ")+1 : strings.LastIndex(adsTxt, " ")])
		lastPage /= psize
	}
	links := make([][]string, lastPage)
	links[0] = page

	wg := &sync.WaitGroup{}
	for i := 2; i <= lastPage; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			links[i-1], _ = listParse(url, i)
		}(i)
	}
	wg.Wait()

	ads := make([][]*domain.Row, lastPage)
	for i := range ads {
		ads[i] = make([]*domain.Row, len(links[i]))
	}
	for i, page := range links {
		for j, link := range page {
			wg.Add(1)
			go func(i, j int, link string) {
				defer wg.Done()
				ads[i][j] = Ad(lat, lng, segment, link, room)
			}(i, j, link)
		}
	}
	wg.Wait()
	result := make([]*domain.Row, 0)
	for _, adSlice := range ads {
		for _, ad := range adSlice {
			if ad != nil {
				result = append(result, ad)
			}
		}
	}
	return result
}

func listParse(url string, page int) (links []string, doc *html.Node) {
	links = make([]string, 0)
	doc, err := htmlquery.LoadURL(url + fmt.Sprintf(params, page))
	if err != nil {
		log.Error("Bad url:", err)
		return
	}

	list := htmlquery.QuerySelectorAll(doc, divLinkArea)
	for _, el := range list {
		links = append(links, htmlquery.SelectAttr(el, "href"))
	}
	return
}

func Ad(lat, lng float64, segment, url string, rooms int) *domain.Row {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Error("Bad url:", err)
		return nil
	}
	priceNode := htmlquery.SelectAttr(htmlquery.QuerySelector(doc, spanPrice), "content")
	priceIdx := strings.LastIndex(priceNode, " ")
	if priceIdx < 0 {
		priceIdx = len(priceNode)
	}
	price, _ := strconv.Atoi(strings.ReplaceAll(priceNode[:priceIdx], " ", ""))

	geoNode := htmlquery.QuerySelector(doc, divGeo)
	if geoNode == nil {
		return nil
	}
	metroNode := htmlquery.QuerySelector(geoNode, li1)
	if metroNode == nil {
		return nil
	}
	metroStr := htmlquery.InnerText(metroNode)
	metroIdx := strings.Index(metroStr, bywalk)
	if metroIdx < 0 {
		return nil
	}
	metroStr = strings.FieldsFunc(metroStr[metroIdx-minutelen:metroIdx], func(r rune) bool {
		return !unicode.IsNumber(r)
	})[0]
	metro, _ := strconv.Atoi(metroStr)
	addr := htmlquery.SelectAttr(htmlquery.QuerySelector(geoNode, spanName), "content")

	resp, err := http.Get(fmt.Sprintf(config.New().CoordApi, price, price, rooms))
	if err != nil {
		log.Error(err)
		return nil
	}
	defer resp.Body.Close()
	coords := new(domain.Coordinates)
	json.NewDecoder(resp.Body).Decode(coords)

	result := &domain.Row{
		Address: addr,
		Cost:    price,
		Balcony: domain.No,
		Metro:   float64(metro),
	}
	for k, v := range coords.Data.Points {
		if strings.Contains(addr, v.Content.Address) {
			i := strings.Index(k, " ")
			result.Latitude, _ = strconv.ParseFloat(k[:i], 64)
			result.Longitude, _ = strconv.ParseFloat(k[i+1:], 64)
		}
	}
	if geo.NewPoint(lat, lng).GreatCircleDistance(geo.NewPoint(result.Latitude, result.Longitude)) > 4 {
		return nil
	}

	summary := htmlquery.QuerySelectorAll(doc, divInfo)
	for _, el := range summary {
		title := htmlquery.QuerySelector(el, divTitle)
		value := htmlquery.QuerySelector(el, divValue)
		switch htmlquery.InnerText(title) {
		case "Общая":
			s := htmlquery.InnerText(value)
			s = s[:strings.Index(s, unispace)]
			result.Total, _ = strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
		case "Кухня":
			s := htmlquery.InnerText(value)
			s = s[:strings.Index(s, unispace)]
			result.Kitchen, _ = strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
		case "Этаж":
			floor := htmlquery.InnerText(value)
			cfloor, _ := strconv.Atoi(floor[:strings.Index(floor, " ")])
			result.CFloor = uint(cfloor)
		case "Построен":
			year, _ := strconv.Atoi(htmlquery.InnerText(value))
			switch {
			case year >= time.Now().Year()-2:
				result.Segment = domain.SegmentNew
			case year <= yearOld:
				result.Segment = domain.SegmentOld
			default:
				result.Segment = domain.SegmentMid
			}
		}
	}
	addFeature := htmlquery.QuerySelectorAll(doc, liFeature)
	for _, el := range addFeature {
		title := htmlquery.QuerySelector(el, span1)
		value := htmlquery.QuerySelector(el, span2)
		switch htmlquery.InnerText(title) {
		case "Балкон/лоджия":
			result.Balcony = domain.Yes
		case "Ремонт":
			switch strings.ToLower(htmlquery.InnerText(value)) {
			case "евроремонт", "дизайнерский":
				result.State = domain.StateNew
			case "косметический":
				result.State = domain.StateMun
			case "без ремонта":
				result.State = domain.StateOff
			default:
				result.State = domain.StateMun
			}
		}
	}
	if result.State == "" {
		result.State = domain.StateMun
	}
	items := htmlquery.QuerySelectorAll(doc, divItem)
	for _, el := range items {
		title := htmlquery.QuerySelector(el, div1)
		value := htmlquery.QuerySelector(el, div2)
		switch htmlquery.InnerText(title) {
		case "Год постройки":
			year, _ := strconv.Atoi(htmlquery.InnerText(value))
			switch {
			case year >= time.Now().Year()-2:
				result.Segment = domain.SegmentNew
			case year <= yearOld:
				result.Segment = domain.SegmentOld
			default:
				result.Segment = domain.SegmentMid
			}
		}
	}
	if result.Segment != segment {
		return nil
	}
	return result
}

var (
	segments = map[string]int{
		domain.LowerSegmentNew: 2,
		domain.LowerSegmentMid: 1,
		domain.LowerSegmentOld: 1,
	}
	walls = map[string]int{
		domain.LowerWallBrick: 1,
		domain.LowerWallMono:  2,
		domain.LowerWallPanel: 3,
	}
	roomType = map[string]int{
		domain.LowerStudio: 9,
	}
)

func FindAnalogs(pattern *domain.Row) []*domain.Row {
	house := walls[strings.ToLower(pattern.Walls)]
	object := segments[strings.ToLower(pattern.Segment)]
	room, err := strconv.Atoi(pattern.Rooms)
	if err != nil {
		room = roomType[pattern.Rooms]
	}
	url := fmt.Sprintf(config.New().SearchApi, house, pattern.Floors, pattern.Floors, object, room)
	return Search(pattern.Latitude, pattern.Longitude, pattern.Segment, url, room)
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

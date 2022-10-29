package htmlparser

import (
	"fmt"
	"leaders_apartments/internal/pkg/domain"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/html"
)

const params = "?page=%d&subcategory=1&limit=%s&sort=13"

func Search(url string) *domain.SearchPage {
	links := make([]string, 0)
	lastPage, limitPage := 1, "30"
	doc := ListParse(&links, url, limitPage, lastPage)
	if doc == nil {
		return nil
	}
	expr, _ := xpath.Compile(`count(//a[@class="pagination-block__page-link"])`)
	v := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(float64)
	page, _ := htmlquery.Query(doc, fmt.Sprintf(`//a[@class="pagination-block__page-link"][%f]`, v))
	lastPage, _ = strconv.Atoi(htmlquery.SelectAttr(page, "data-page"))
	limitPage = htmlquery.SelectAttr(page, "data-limit")
	for i := 2; i <= lastPage; i++ {
		if err := ListParse(&links, url, limitPage, i); err == nil {
			return nil
		}
	}
	return &domain.SearchPage{Links: links, Count: len(links)}
}

func ListParse(links *[]string, url, limit string, page int) *html.Node {
	doc, err := htmlquery.LoadURL(url + fmt.Sprintf(params, page, limit))
	if err != nil {
		log.Error("Bad url:", err)
		return nil
	}
	if list, err := htmlquery.QueryAll(doc, `//a[@class="search-item__item-link"]`); err == nil {
		for _, el := range list {
			*links = append(*links, "https:"+htmlquery.SelectAttr(el, "href"))
		}
	} else {
		log.Error(err)
		return nil
	}
	return doc
}

func Ad(url, floor string) *domain.AdPage {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Error("Bad url:", err)
		return nil
	}
	if coord, err := htmlquery.Query(doc, `//div[@class="yamap-container"]`); err == nil {
		data := htmlquery.SelectAttr(coord, "data-thumb")
		begin := strings.Index(data, "=") + 1
		end := begin + strings.Index(data[begin:], "&")
		coords := strings.Split(data[begin:end], ",")
		return &domain.AdPage{Longitude: coords[0], Latitude: coords[1]}
	} else {
		log.Error(err)
	}
	return nil
}

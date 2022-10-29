package htmlparser

import (
	"leaders_apartments/internal/pkg/domain"
	"net/http"

	"github.com/antchfx/htmlquery"
	"github.com/labstack/gommon/log"
)

func Search(url string) *domain.SearchPage {
	links := make([]string, 0)
	lastPage := ""
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Error("Bad url:", err)
		return nil
	}
	list, err := htmlquery.QueryAll(doc, "//a[@class = search-item__title-link search-item__item-link]")
	for _, el := range list {
		links = append(links, "https:"+htmlquery.SelectAttr(el, "href"))
	}
	log.Info(lastPage)
	return &domain.SearchPage{Links: links}
}

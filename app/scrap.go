// scrap.go
package app

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Scrap() (float64, error) {
	url := "https://alsuper.com/producto/pulpa-de-res-picada-357825"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, fmt.Errorf("status code error: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, err
	}

	var priceText string
	doc.Find(".as-font-24").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "$") {
			priceText = text
			return false // break loop
		}
		return true
	})

	// Extract number from price string, e.g. "$139.90" → 139.90
	re := regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`)
	match := re.FindString(priceText)
	if match == "" {
		return 0, fmt.Errorf("no price found")
	}

	price, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

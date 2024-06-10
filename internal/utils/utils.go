package utils

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"
)

type Article struct {
	Title         string   `json:"title"`
	Url           string   `json:"url"`
	Authors       []string `json:"authors"`
	Description   string   `json:"description"`
	FullText      string   `json:"full_text"`
	PublishedDate string   `json:"published_date"`
	Image         string   `json:"image"`
	Keywords      []string `json:"keywords"`
	RawHTML       string   `json:"raw_html"`
}

func IsValidDate(date_str string) bool {
	match, err := regexp.MatchString("\\d{4}-\\d{1,2}-\\d{1,2}", date_str)
	if err != nil {
		return false
	}
	return match
}

func ScrapeContent(url string, timeoutSeconds time.Duration) *http.Response {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	j, _ := cookiejar.New(nil)
	client := &http.Client{Jar: j, Timeout: timeoutSeconds}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Failed to scrape %s", url)
	}

	if response.StatusCode != 200 {
		log.Fatalf("HTTP Error scraping %s got status code %d", url, response.StatusCode)
	}
	return response
}

func StripTrailing(title string) string {
	separators := []string{" - ", " | "}
	for _, sep := range separators {
		if strings.Contains(title, sep) {
			tail := title[strings.LastIndex(title, sep)+1:]
			orgName := strings.TrimSpace(strings.Trim(tail, sep))
			if len(strings.Split(orgName, " ")) < 5 {
				return title[:strings.LastIndex(title, sep)]
			}
		}
	}
	return title
}

func ExtractLDJson(document *html.Node) (map[string]interface{}, bool) {
	// https://developers.google.com/search/docs/appearance/structured-data/article#json-ld
	jsonLDs := htmlquery.Find(document, "//script[@type=\"application/ld+json\"]")

	var result interface{}
	var results []interface{}

	for _, jsonLD := range jsonLDs {
		if jsonLD == nil {
			continue
		}

		err := json.Unmarshal([]byte(htmlquery.InnerText(jsonLD)), &result)
		if err != nil {
			log.Printf("Error parsing JSON: %v\n", err)
		}

		// convert single object {...} to a list of single object [{...}]
		if obj, ok := result.(map[string]interface{}); ok {
			results = append(results, obj)
		} else if list, ok := result.([]interface{}); ok {
			results = list
		}

		for _, item := range results {
			if obj, ok := item.(map[string]interface{}); ok {
				if jsonType, ok := obj["@type"].(string); ok {
					if jsonType == "NewsArticle" {
						return obj, true
					}
				} else {
					log.Println("@type not found or not a string")
				}
			}
		}
	}

	return nil, false
}

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

func IsValidDate(dateStr string) bool {
	match, err := regexp.MatchString("\\d{4}-\\d{1,2}-\\d{1,2}", dateStr)
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

func RemoveDuplicates(values []string) []string {
	// remove case-insensitive duplicates, blanks
	seen := make(map[string]bool)
	var result []string

	for _, val := range values {
		cleanedAuthor := strings.ToLower(strings.TrimSpace(val))
		if cleanedAuthor == "" {
			continue
		}

		if !seen[cleanedAuthor] {
			seen[cleanedAuthor] = true
			result = append(result, strings.TrimSpace(val))
		}
	}
	return result
}

func RemoveStopWords(words []string) []string {
	var stopWords = []string{
		"a", "about", "above", "after", "again", "against", "all", "am", "an", "and", "any", "are", "aren't", "as",
		"at", "be", "because", "been", "before", "being", "below", "between", "both", "but", "by", "can't", "cannot",
		"could", "couldn't", "did", "didn't", "do", "does", "doesn't", "doing", "don't", "down", "during", "each",
		"few", "for", "from", "further", "had", "hadn't", "has", "hasn't", "have", "haven't", "having", "he", "he'd",
		"he'll", "he's", "her", "here", "here's", "hers", "herself", "him", "himself", "his", "how", "how's", "i",
		"i'd", "i'll", "i'm", "i've", "if", "in", "into", "is", "isn't", "it", "it's", "its", "itself", "let's", "me",
		"more", "most", "mustn't", "my", "myself", "no", "nor", "not", "of", "off", "on", "once", "only", "or",
		"other", "ought", "our", "ours", "ourselves", "out", "over", "own", "same", "shan't", "she", "she'd",
		"she'll", "she's", "should", "shouldn't", "so", "some", "such", "than", "that", "that's", "the", "their",
		"theirs", "them", "themselves", "then", "there", "there's", "these", "they", "they'd", "they'll", "they're",
		"they've", "this", "those", "through", "to", "too", "under", "until", "up", "very", "was", "wasn't", "we",
		"we'd", "we'll", "we're", "we've", "were", "weren't", "what", "what's", "when", "when's", "where", "where's",
		"which", "while", "who", "who's", "whom", "why", "why's", "with", "won't", "would", "wouldn't", "you", "you'd",
		"you'll", "you're", "you've", "your", "yours", "yourself", "yourselves",
	}

	// Convert stop words to a set for efficient look-up
	stopWordsSet := make(map[string]struct{}, len(stopWords))
	for _, word := range stopWords {
		stopWordsSet[strings.ToLower(word)] = struct{}{}
	}

	// Declare an empty slice for filtered keywords
	var filteredKeywords []string

	// Iterate over the keywords and append non-stop words
	for _, word := range words {
		if _, found := stopWordsSet[strings.ToLower(word)]; !found {
			filteredKeywords = append(filteredKeywords, word)
		}
	}
	return filteredKeywords
}

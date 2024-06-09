package extractors

import (
	"github.com/antchfx/htmlquery"
	"github.com/sagunsh/gophernews/internal/utils"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

func ExtractTitle(document *html.Node, jsonLD map[string]interface{}) string {
	titleTag := htmlquery.FindOne(document, "//title")
	metaProperty := htmlquery.FindOne(document, "//meta[@property=\"og:title\"]/@content")
	metaName := htmlquery.FindOne(document, "//meta[@name=\"title\"]/@content")
	h1 := htmlquery.FindOne(document, "//h1/text()")

	titleText := ""
	if titleTag != nil {
		titleText = strings.TrimSpace(htmlquery.InnerText(titleTag))
	}

	metaPropertyText := ""
	if metaProperty != nil {
		metaPropertyText = strings.TrimSpace(htmlquery.InnerText(metaProperty))
	}

	meteNameText := ""
	if metaName != nil {
		meteNameText = strings.TrimSpace(htmlquery.InnerText(metaName))
	}

	h1Text := ""
	if h1 != nil {
		h1Text = strings.TrimSpace(htmlquery.InnerText(h1))
	}

	headline := ""
	if jsonLD != nil {
		if val, ok := jsonLD["headline"].(string); ok {
			headline = strings.TrimSpace(val)
		}
	}

	if h1Text != "" {
		if titleText != "" && strings.Contains(strings.ToLower(titleText), strings.ToLower(h1Text)) {
			return h1Text
		}

		if metaPropertyText != "" && strings.Contains(strings.ToLower(metaPropertyText), strings.ToLower(h1Text)) {
			return h1Text
		}

		if headline != "" && strings.Contains(strings.ToLower(headline), strings.ToLower(h1Text)) {
			return h1Text
		}
	}

	if titleText != "" {
		return titleText
	}

	if metaPropertyText != "" {
		return metaPropertyText
	}

	if meteNameText != "" {
		return meteNameText
	}

	return ""
}

func ExtractAuthors(document *html.Node, jsonLD map[string]interface{}) []string {
	var authorList []string

	if jsonLD != nil {
		if authors, ok := jsonLD["author"].([]interface{}); ok {
			for _, author := range authors {
				if authorInfo, ok := author.(map[string]interface{}); ok {
					if authorName, ok := authorInfo["name"].(string); ok {
						authorList = append(authorList, authorName)
					}
				}
			}
		}
	}

	if len(authorList) > 0 {
		return authorList
	}

	return authorList
}

func ExtractDescription(document *html.Node, jsonLD map[string]interface{}) string {
	return ""
}

func ExtractFullText(document *html.Node, jsonLD map[string]interface{}) string {
	return ""
}

func ExtractPublishedDate(document *html.Node, jsonLD map[string]interface{}) string {
	// check in meta tag, <time datetime"..."> and ld+json datePublished
	dateXpaths := []string{"//meta[@property='article:published_time']/@content", "//time/@datetime"}
	for _, xpath := range dateXpaths {
		date := htmlquery.FindOne(document, xpath)
		if date != nil {
			if utils.IsValidDate(htmlquery.InnerText(date)) {
				return htmlquery.InnerText(date)
			}
		}
	}

	if jsonLD != nil {
		if dateStr, ok := jsonLD["datePublished"].(string); ok {
			if utils.IsValidDate(dateStr) {
				return dateStr
			}
		}
	}

	return ""
}

func ExtractImage(document *html.Node, jsonLD map[string]interface{}) string {
	return ""
}

func ExtractKeywords(document *html.Node, jsonLD map[string]interface{}) []string {
	var keywords []string
	return keywords
}

func ExtractRawHTML(response *http.Response) string {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

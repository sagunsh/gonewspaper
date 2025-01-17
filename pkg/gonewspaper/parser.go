package gonewspaper

import (
	"bytes"
	"github.com/antchfx/htmlquery"
	"github.com/sagunsh/gonewspaper/internal/extractors"
	"github.com/sagunsh/gonewspaper/internal/utils"
	"io"
	"log"
	"strings"
	"time"
)

func ParseArticle(url string, timeoutSeconds time.Duration) utils.Article {
	response := utils.ScrapeContent(url, timeoutSeconds)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	rawHTML := ""
	if err == nil {
		rawHTML = strings.TrimSpace(string(body))
	}

	document, err := htmlquery.Parse(bytes.NewReader([]byte(rawHTML)))
	if err != nil {
		log.Fatal(err)
	}

	ldJson, _ := utils.ExtractLDJson(document)

	article := utils.Article{}
	article.Title = utils.StripTrailing(extractors.ExtractTitle(document, ldJson))
	article.Url = response.Request.URL.String()
	article.Authors = extractors.ExtractAuthors(document, ldJson)
	article.Description = extractors.ExtractDescription(document, ldJson)
	article.FullText = extractors.ExtractFullText(document, ldJson)
	article.PublishedDate = extractors.ExtractPublishedDate(document, ldJson)
	article.Image = extractors.ExtractImage(document, ldJson)
	article.Keywords = extractors.ExtractKeywords(document, ldJson)
	article.RawHTML = rawHTML[0:100]
	return article
}

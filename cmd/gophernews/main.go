package main

import (
	"encoding/json"
	"fmt"
	"github.com/sagunsh/gophernews/pkg/gophernews"
	"log"
	"os"
)

func main() {
	//"https://kathmandupost.com/province-no-2/2024/06/04/unified-socialist-quits-madhesh-government"
	//"https://www.abc.net.au/news/2024-06-03/josh-frydenberg-canberra-comeback-kooyong-amelia-hamer/103928586"
	//"https://www.skynews.com.au/australia-news/young-australians-shocked-and-outraged-at-video-of-university-students-in-the-1970s-revealing-their-living-expenses/news-story/6726ce45c45198c02594c994b9a93585"
	//"https://edition.cnn.com/2024/06/05/india/india-general-election-result-modi-analysis-intl-hnk/index.html"

	var url string
	if len(os.Args) > 1 {
		url = os.Args[1]
	} else {
		fmt.Print("Enter a news URL: ")
		_, err := fmt.Scan(&url)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	article := gophernews.ParseArticle(url)
	// fmt.Println(article)
	jsonData, err := json.MarshalIndent(article, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

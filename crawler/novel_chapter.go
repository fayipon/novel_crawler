package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// 要抓取的網址
	url := "https://sto520.com/book/27711/"

	// 發送 HTTP GET 請求
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	
	fmt.Printf(" %s\n", response)

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 找到所有的 <dd> 元素並取出文字
	doc.Find("dd").Each(func(index int, element *goquery.Selection) {
		text := element.Text()
		fmt.Printf("Text in <dd> %d: %s\n", index+1, text)
	})
}
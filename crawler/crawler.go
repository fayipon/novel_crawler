package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/PuerkitoBio/goquery"
)

func main() {
    url := "https://www.85novel.com/book/2542416.html"

    // 发送HTTP GET请求
    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()

    // 使用goquery解析HTML
    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        log.Fatal(err)
    }

    // 找到所有包含/book/2542416/的链接
    doc.Find("a[href^='/book/2542416/']").Each(func(index int, item *goquery.Selection) {
        // 获取链接的href属性
        href, _ := item.Attr("href")
        // 获取链接的title文本
        title := item.Text()
        // 输出链接的href和title
        fmt.Printf("链接: %s\n标题: %s\n", href, title)
    })
}
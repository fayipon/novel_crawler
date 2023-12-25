package main

import (
    "fmt"
    "log"
    "github.com/PuerkitoBio/goquery"
    "net/http"
)

func main() {
    // 要抓取的网站URL
    url := "https://www.bg3.co/novel/pagea/chaoshenjixieshi-qipeijia_1.html"

    // 发起HTTP GET请求
    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()

    // 检查响应状态码
    if res.StatusCode != http.StatusOK {
        log.Fatalf("HTTP request failed with status code: %d", res.StatusCode)
    }

    // 使用goquery选择器来解析页面元素
    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        log.Fatal(err)
    }

    // 选择所有带有<p style="color: black;">的元素
    doc.Find("p[style='color: black;']").Each(func(index int, element *goquery.Selection) {
        // 提取文本内容并打印
        text := element.Text()
        fmt.Println("文本内容:", text)
    })
}
package main

import (
    "fmt"
    "log"
    "github.com/PuerkitoBio/goquery"
    "net/http"
)

func main() {
    // 要抓取的网站URL
    url := "https://sto520.com/book/27711/"

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

    // 使用goquery选择器来解析页面元素
    // 这里以获取标题为例
    title := doc.Find("h1").Text()

    // 输出标题
    fmt.Println("标题:", title)
}
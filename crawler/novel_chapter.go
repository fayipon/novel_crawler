package main

import (
    "fmt"
    "log"
    "github.com/PuerkitoBio/goquery"
)

func main() {
    // 要抓取的网站URL
    url := "https://sto520.com/book/27711/"

    // 发起HTTP GET请求
    doc, err := goquery.NewDocument(url)
    if err != nil {
        log.Fatal(err)
    }

    // 使用goquery选择器来解析页面元素
    // 这里以获取标题为例
    title := doc.Find("h1").Text()
    
    // 输出标题
    fmt.Println("标题:", title)
}
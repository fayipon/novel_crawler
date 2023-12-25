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

    // 创建自定义请求
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    // 添加HTTP头部信息
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Accept-Encoding", "gzip, deflate, br")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
    req.Header.Set("Host", "sto520.com")
    req.Header.Set("Cookie", "__cf_bm=YEtKER0mn7qiQdBXPEBX.MmkCarTbagRi3I0ZbVNtjE-1703471246-1-AWY2Z+ZJY2VbDywXFt10jMiYdff6K/GlXFiXWHOSKBfsTgOnBl1zPL6Xx5x/Pq13FSK1PpnIFZkArHZLDvscgxI=")

    // 发起HTTP请求
    client := &http.Client{}
    res, err := client.Do(req)
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
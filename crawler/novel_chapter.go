package main

import (
    "fmt"
    "log"
    "os"
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

    // 获取页面标题
    title := doc.Find("title").Text()

    // 指定目录路径
    dirPath := "./novel" // 例如，将文件夹保存在当前工作目录下的 "novel" 目录中

    // 创建目录（如果不存在）
    err = os.MkdirAll(dirPath, os.ModePerm)
    if err != nil {
        log.Fatal(err)
    }

    // 创建txt文件并打开以供写入
    filePath := fmt.Sprintf("%s/%s.txt", dirPath, title)
    file, err := os.Create(filePath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // 选择所有带有<p style="color: black;">的元素
    doc.Find("p").Each(func(index int, element *goquery.Selection) {
        // 提取文本内容
        text := element.Text()
        
        // 写入文本内容到txt文件中
        _, err := file.WriteString(text + "\n")
        if err != nil {
            log.Fatal(err)
        }
    })

    fmt.Printf("文本已保存到 %s\n", filePath)
}
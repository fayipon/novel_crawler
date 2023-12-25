package main

import (
    "fmt"
    "log"
    "os"
    "github.com/PuerkitoBio/goquery"
    "net/http"
    "strconv"
)

func main() {
    // 指定基本URL
    baseURL := "https://www.bg3.co/novel/pagea/chaoshenjixieshi-qipeijia_"

    // 循环从1到5
    for i := 1; i <= 5; i++ {
        // 构建完整URL
        url := baseURL + strconv.Itoa(i) + ".html"

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

        // 选择您希望提取的内容并进行处理
        // 这里可以根据您的需求提取文本内容并进行处理

        // 示例：获取页面标题并打印
        title := doc.Find("title").Text()
        fmt.Printf("标题 %d: %s\n", i, title)

        // 将页面内容保存为txt文件
        filePath := fmt.Sprintf("page%d.txt", i)
        file, err := os.Create(filePath)
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()

        // 提取文本内容并写入txt文件
        doc.Find("p").Each(func(index int, element *goquery.Selection) {
            text := element.Text()
            _, err := file.WriteString(text + "\n")
            if err != nil {
                log.Fatal(err)
            }
        })

        fmt.Printf("页面内容已保存到 %s\n", filePath)
    }
}
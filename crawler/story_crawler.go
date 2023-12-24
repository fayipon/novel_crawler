package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/PuerkitoBio/goquery"
    "regexp"
    "strings"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
    "os"
	"strconv"
)

const siteID = 2542416 // 如果 site_id 是常量，可以定义为常量

func main() {
    // 加载 .env 文件中的配置信息
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // 获取从 .env 文件加载的配置信息
    username := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME") // 添加数据库名称配置

    // MySQL 数据库连接信息
    dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", username, password, dbname)

    // 连接到 MySQL 数据库
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	_, err = db.Exec("USE novel")
	if err != nil {
		log.Fatal(err)
	}

    siteURL := "https://www.85novel.com/book/2542416.html"

    // 发送HTTP GET请求
    res, err := http.Get(siteURL)
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

        // 使用正则表达式提取数字
        re := regexp.MustCompile(`(\d+)`)
        matches := re.FindAllString(href, -1)
        storyID := 0
        if len(matches) > 0 {
            // 提取的数字
            storyIDStr := matches[len(matches)-1]
            storyID, _ = strconv.Atoi(storyIDStr)
        }

        // 使用空格分割标题
        titleParts := strings.Fields(title)
        if len(titleParts) > 1 {
            chapterName := titleParts[0]
            storyName := titleParts[1]

            // 将数字和标题写入 MySQL 数据库
            insertSQL := "INSERT INTO story (site_id, story_id, chapter_name, story_name) VALUES (?, ?, ?, ?)"
            _, err := db.Exec(insertSQL, siteID, storyID, chapterName, storyName)
            if err != nil {
                log.Fatal(err)
            }
            fmt.Printf("数据已成功写入 MySQL 数据库 - 章节: %s, 标题: %s\n", chapterName, storyName)
        }
    })
}
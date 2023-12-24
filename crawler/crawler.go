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
)

func main() {

	// 加载 .env 文件中的配置信息
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 获取从 .env 文件加载的配置信息
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	// MySQL 数据库连接信息
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/dbname", username, password)

	// 连接到MySQL数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	///////////////////////

	site_id := 2542416
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
       // fmt.Printf("链接: %s\n标题: %s\n", href, title)
		
		// 使用正则表达式提取数字
		re := regexp.MustCompile(`(\d+)`)
		matches := re.FindAllString(href, -1)
		story_id := 0
		if len(matches) > 0 {
			// 提取的数字
			story_id = matches[len(matches)-1]
			fmt.Printf("数字: %s\n", story_id)
		}

		// 使用空格分割标题
		titleParts := strings.Fields(title)
		if len(titleParts) > 0 {
			// 分割的标题部分
			fmt.Printf("章節: %v\n", titleParts[0])
			fmt.Printf("標題: %v\n", titleParts[1])

			// 将数字和标题写入 MySQL 数据库
			insertSQL := "INSERT INTO story (site_id, story_id, story_name, story_url) VALUES (?, ?, ?, ?)"
			_, err := db.Exec(insertSQL, site_id, story_id, titleParts[0], titleParts[1])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("数据已成功写入 MySQL 数据库")
		}

    })
}
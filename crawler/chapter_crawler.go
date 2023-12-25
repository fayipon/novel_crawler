package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
    "os"
    "time"
    "github.com/PuerkitoBio/goquery"
)

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

    // 建立数据库连接
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    // 执行查询操作
    rows, err := db.Query("SELECT * FROM story WHERE status = 1 limit 1")
    if err != nil {
        panic(err.Error())
    }
    defer rows.Close()

    // 处理查询结果
    for rows.Next() {
        var id int
        var site_id int
        var story_id int
        var story_name string
        var chapter_name string // Change "title" to "chapter_name"
        var status int
        // 添加其他需要的字段

        err := rows.Scan(&id, &site_id, &story_id, &story_name, &chapter_name, &status) // Fix variable names
        if err != nil {
            panic(err.Error())
        }

        // 处理查询结果，可以根据需要输出或进行其他操作
        fmt.Printf("ID: %d, story_name: %s, chapter_name: %s\n", id, story_name, chapter_name)

        link := fmt.Sprintf("https://www.85novel.com/book/%d/%d.html", site_id, story_id)
        fmt.Printf("URL: %s\n", link)

        // 发送 HTTP 请求获取页面内容
        doc, err := goquery.NewDocument(link)
        if err != nil {
            log.Fatal(err)
        }

        // 创建一个字符串变量来保存抓取到的文本
        var mergedText string

        // 查找所有的 <p></p> 标签并提取文本
        doc.Find("p").Not("p:last-child").Each(func(index int, element *goquery.Selection) {
            text := element.Text()
            mergedText += text + "\n" 
        })

        // 连接数据库
        db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/dbname")
        if err != nil {
            log.Fatal(err)
        }
        defer db.Close()

        // 插入数据到数据库表
        currentTime := time.Now().Format("2006-01-02 15:04:05")
        _, err = db.Exec("INSERT INTO your_table (site_id, story_id, data, create_time) VALUES (?, ?, ?, ?)", site_id, story_id, mergedText, currentTime)
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Println("ID: %d, story_name: %s, chapter_name: %s \n Data inserted successfully!\n", id, story_name, chapter_name)

    }

    // 检查是否有错误发生
    if err := rows.Err(); err != nil {
        panic(err.Error())
    }
}
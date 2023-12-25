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
        doc.Find("p").Each(func(index int, element *goquery.Selection) {
            text := element.Text()
            mergedText += text + "\n" 
        })

        // 移除指定文本
        textToRemove := "85小說網2013- All Rights Reserved 免責聲明: 本站所有文學作品均由自動化程式以非人工方式自動從第三方網頁搜索而成，所有文學作品都可以在第三方網頁上以公開的方式被找到， 本站收錄這些文學作品是為了便於廣大書友交流學習，不代表本站同意這些文學作品的立場和內容。 任何單位或個人認為本站收錄到的第三方網頁內容可能涉嫌侵犯其信息網路傳播權，應該及時向本站提出書面權利通知，並提供身份證明、權屬證明及詳細侵權情況證明。本站在收到上述法律文件後，將會依法儘快斷開相關鏈接內容。 本站尊重並保護所有使用本站用戶的個人隱私權，您註冊的用戶名、電子郵件地址等個人資料，非經您親自許可或根據相關法律、法規的強制性規定，本站不會主動地泄露給第三方。 本站所收錄小說作品、社區話題、小說評論及本站所做之廣告均屬其個人行為，與本站立場無關"
        mergedText = strings.Replace(mergedText, textToRemove, "", -1)

        fmt.Printf("%s\n", mergedText)

        // 插入数据到数据库表
        currentTime := time.Now().Format("2006-01-02 15:04:05")
        _, err = db.Exec("INSERT INTO chapter (site_id, story_id, data, create_time) VALUES (?, ?, ?, ?)", site_id, story_id, mergedText, currentTime)
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Println("ID: %d, story_name: %s, chapter_name: %s \n Data inserted successfully!\n", id, story_name, chapter_name)

        // 执行 SQL 更新操作
        _, err = db.Exec("UPDATE story SET status = '2' WHERE site_id = ? AND story_id = ?", site_id, story_id,)
        if err != nil {
            log.Fatal(err)
        }

    }

    // 检查是否有错误发生
    if err := rows.Err(); err != nil {
        panic(err.Error())
    }
}
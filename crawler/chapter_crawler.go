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

// Story 结构体用于表示一个故事
type Story struct {
    ID         int
    SiteID     int
    StoryID    int
    StoryName  string
    ChapterName string
    Status     int
}

// Chapter 结构体用于表示一个章节
type Chapter struct {
    SiteID    int
    StoryID   int
    Data      string
    CreateTime string
}

func main() {
    // 加载 .env 文件中的配置信息
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // 获取从 .env 文件加载的配置信息
    username := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    // MySQL 数据库连接信息
    dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", username, password, dbname)

    // 建立数据库连接
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 获取待处理的故事列表
    stories, err := getStoriesToProcess(db)
    if err != nil {
        log.Fatal(err)
    }

    // 处理每个故事
    for _, story := range stories {
        processStory(db, story)
    }
}

// getStoriesToProcess 从数据库中获取待处理的故事列表
func getStoriesToProcess(db *sql.DB) ([]Story, error) {
    rows, err := db.Query("SELECT * FROM story WHERE status = 1 LIMIT 1")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stories []Story
    for rows.Next() {
        var story Story
        if err := rows.Scan(&story.ID, &story.SiteID, &story.StoryID, &story.StoryName, &story.ChapterName, &story.Status); err != nil {
            return nil, err
        }
        stories = append(stories, story)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return stories, nil
}

// processStory 处理故事
func processStory(db *sql.DB, story Story) {
    fmt.Printf("Processing Story ID: %d, Story Name: %s, Chapter Name: %s\n", story.ID, story.StoryName, story.ChapterName)

    link := fmt.Sprintf("https://www.85novel.com/book/%d/%d.html", story.SiteID, story.StoryID)
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

    specificString := " 85小說網2013- All Rights Reserved  免責聲明: 本站所有文學作品均由自動化程式以非人工方式自動從第三方網頁搜索而成，所有文學作品都可以在第三方網頁上以公開的方式被找到， 本站收錄這些文學作品是為了便於廣大書友交流學習，不代表本站同意這些文學作品的立場和內容。 任何單位或個人認為本站收錄到的第三方網頁內容可能涉嫌侵犯其信息網路傳播權，應該及時向本站提出書面權利通知，並提供身份證明、權屬證明及詳細侵權情況證明。本站在收到上述法律文件後，將會依法儘快斷開相關鏈接內容。 本站尊重並保護所有使用本站用戶的個人隱私權，您註冊的用戶名、電子郵件地址等個人資料，非經您親自許可或根據相關法律、法規的強制性規定，本站不會主動地泄露給第三方。 本站所收錄小說作品、社區話題、小說評論及本站所做之廣告均屬其個人行為，與本站立場無關。 "

    // 检查 mergedText 是否等于特定字符串
    if mergedText == specificString {
        // 在 mergedText 等于特定字符串时执行其他分析逻辑
        fmt.Println("Merged text matches the specific string. Performing additional analysis...")

        // 在此处添加其他分析逻辑
      
    } else {

        // 插入数据到数据库表
        currentTime := time.Now().Format("2006-01-02 15:04:05")
        _, err = db.Exec("INSERT INTO chapter (site_id, story_id, data, create_time) VALUES (?, ?, ?, ?)", story.SiteID, story.StoryID, mergedText, currentTime)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Data inserted successfully!\n")

        // 执行 SQL 更新操作
        _, err = db.Exec("UPDATE story SET status = 2 WHERE site_id = ? AND story_id = ?", story.SiteID, story.StoryID)
        if err != nil {
            log.Fatal(err)
        }
    }

}
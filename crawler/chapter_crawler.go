package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
    "os"
    "time"
    "strings"
    "regexp"
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
    var contentInfoJSON string // 声明 contentInfoJSON 变量

    // 查找所有的 <p></p> 标签并提取文本
    doc.Find("p").Each(func(index int, element *goquery.Selection) {
        text := element.Text()
        mergedText += text + "\n"
    })

    specificPrefix := " 85小說網"


    // 检查 mergedText 是否以特定前缀开始
    if strings.HasPrefix(mergedText, specificPrefix) {
        // 在 mergedText 以特定前缀开始时执行其他处理逻辑
        fmt.Println("Merged text starts with the specific prefix. Performing additional analysis...")


        // 从 GoQuery 文档中获取包含 JavaScript 代码的 script 标签
        doc.Find("script").Each(func(index int, scriptElement *goquery.Selection) {
            // 获取 script 标签内的文本内容
            scriptText := scriptElement.Text()

            // 使用正则表达式模式，匹配包含 contentInfo 的 JavaScript 代码
            pattern := `contentInfo:\s+(\[[^\]]+\])`

            // 使用正则表达式查找匹配
            re := regexp.MustCompile(pattern)
            match := re.FindStringSubmatch(scriptText)

            if len(match) == 2 {
                // 去掉 JSON 字符串中的空格和换行符
                contentInfoJSON = match[1]
                contentInfoJSON = strings.ReplaceAll(contentInfoJSON, " ", "")
                contentInfoJSON = strings.ReplaceAll(contentInfoJSON, "\n", "")
    
                // 反转字符串
                contentInfoJSONReversed := reverseString(contentInfoJSON)
    
                // 分割 JSON 数组中的字符串，并按行打印
                lines := strings.Split(contentInfoJSONReversed, `","`)
                for _, line := range lines {
                    // 去掉 JSON 数组中的引号和方括号
                    line = strings.Trim(line, `[]"`)
                    // 打印每行
                    fmt.Println(line)
                }
            }
        })

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

// 反转字符串的函数
func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
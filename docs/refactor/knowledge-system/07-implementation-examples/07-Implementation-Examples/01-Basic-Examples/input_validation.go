package main

import (
	"database/sql"
	"fmt"
	"html"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

// 用户名校验（只允许字母、数字、下划线，3-16位）
func isValidUsername(username string) bool {
	pattern := `^[a-zA-Z0-9_]{3,16}$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// 防止SQL注入（使用参数化查询）
func safeQuery(db *sql.DB, username string) (string, error) {
	var email string
	row := db.QueryRow("SELECT email FROM users WHERE username = ?", username)
	if err := row.Scan(&email); err != nil {
		return "", err
	}
	return email, nil
}

// 防止XSS攻击（转义HTML）
func safeHTML(input string) string {
	return html.EscapeString(input)
}

func main() {
	username := "user_01"
	fmt.Println("用户名校验:", isValidUsername(username))

	// 模拟数据库操作
	db, _ := sql.Open("sqlite3", ":memory:")
	db.Exec("CREATE TABLE users (username TEXT, email TEXT)")
	db.Exec("INSERT INTO users VALUES (?, ?)", username, "user01@example.com")

	email, err := safeQuery(db, username)
	if err != nil {
		fmt.Println("查询失败:", err)
	} else {
		fmt.Println("查询到邮箱:", email)
	}

	// XSS防护示例
	unsafeInput := "<script>alert('xss')</script>"
	safe := safeHTML(unsafeInput)
	fmt.Println("原始输入:", unsafeInput)
	fmt.Println("转义后:", safe)
}

package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"

	// 生成密码哈希
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	hash := string(bytes)

	fmt.Println("管理员密码哈希值:")
	fmt.Println(hash)
	fmt.Println()
	fmt.Println("SQL 插入语句:")
	fmt.Printf("INSERT INTO users (username, password, real_name, role, phone, email, status, created_at, updated_at) VALUES ('admin', '%s', '系统管理员', 'supervisor', '13800138000', 'admin@example.com', 1, NOW(), NOW());\n", hash)
}

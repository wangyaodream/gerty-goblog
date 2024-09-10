package main

import (
	"fmt"
	"os"
)

func main() {
	envVars := []string{
		"DB_HOST",
        "DB_NAME",
		"DB_USER",
		"DB_PASSWORD",
		"DB_PORT",
		"DB_TYPE",
	}
	file, err := os.Create(".env")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for _, key := range envVars {
		// 每个环境变量后面跟一个等号和一对空的双引号
		file.WriteString(fmt.Sprintf("%s=\"\"\n", key))
	}

	fmt.Println("file created successfully")
}

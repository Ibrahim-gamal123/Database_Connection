package main

import (
	"bufio"
	"distributed-db/shared"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Starting Master...")

	db := shared.NewDBHandler("root", "", "127.0.0.1:3306")

	go StartServer(db)      // تشغيل TCP server
	go startWebServer(db)   // تشغيل واجهة الويب

	fmt.Println("You can type queries directly into the master:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		query := scanner.Text()
		if query == "exit" {
			break
		}
		req := shared.Request{
			Token:     "secret-token",
			Query:     query,
			FromSlave: "master",
			IsSelect:  strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "SELECT"),
		}
		resp := handleLocalQuery(req, db)
		fmt.Printf("Status: %s | Message: %s\n", resp.Status, resp.Message)
		if resp.Status == "ok" && req.IsSelect {
			fmt.Println(resp.Header)
			for _, row := range resp.Rows {
				fmt.Println(row)
			}
		}
	}
}

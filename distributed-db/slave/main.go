package main

import (
	"bufio"
	
	"fmt"
	"net/http"
	"os"
	
)

func main() {
	// واجهة المستخدم الرسومية
	go func() {
		fs := http.FileServer(http.Dir("slave/web"))
		http.Handle("/", fs)

		// API لتلقي الاستعلامات من الواجهة
		http.HandleFunc("/api/query", handleAPIQuery)

		fmt.Println("Slave GUI running at http://localhost:8081/")
		http.ListenAndServe(":8081", nil)
	}()

	// التشغيل من الطرفية
	fmt.Println("Slave started. Type query:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		query := scanner.Text()
		if query == "exit" {
			break
		}
		sendQueryToMaster(query)
	}
}
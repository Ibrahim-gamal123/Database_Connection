package main

import (
	"bufio"
	"distributed-db/shared"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"net/http"
)

func sendQueryToMaster(query string) string {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		return fmt.Sprintf("Could not connect to master: %v", err)
	}
	defer conn.Close()

	req := shared.Request{
		Token:     "secret-token",
		Query:     query,
		FromSlave: "slave1",
		IsSelect:  strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "SELECT"),
	}
	data, _ := json.Marshal(req)
	fmt.Fprintf(conn, string(data)+"\n")

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Sprintf("Failed to read response from master: %v", err)
	}
	return response
}

// معالج API للواجهة
func handleAPIQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Query string `json:"query"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := sendQueryToMaster(req.Query)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}
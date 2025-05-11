package main

import (
	"bufio"
	"distributed-db/shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const validToken = "secret-token"

func isMasterQuery(query string) bool {
	query = strings.ToUpper(strings.TrimSpace(query))
	return strings.HasPrefix(query, "CREATE") || strings.HasPrefix(query, "DROP")
}

func startWebServer(db *shared.DBHandler) {
	http.Handle("/", http.FileServer(http.Dir("./master/web")))
	http.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req shared.Request
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		req.Token = "secret-token"
		req.FromSlave = "master"
		req.IsSelect = strings.HasPrefix(strings.ToUpper(strings.TrimSpace(req.Query)), "SELECT")
		resp := handleLocalQuery(req, db)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	http.ListenAndServe(":8080", nil)
}

func handleSlave(conn net.Conn, db *shared.DBHandler, logger *os.File) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		var req shared.Request
		err := json.Unmarshal([]byte(line), &req)
		if err != nil || req.Token != validToken {
			log.Println("Rejected request (invalid token or format)")
			conn.Write([]byte("{\"status\":\"error\",\"message\":\"unauthorized\"}\n"))
			continue
		}

		if isMasterQuery(req.Query) && req.FromSlave != "master" {
			log.Println("Rejected request (only master can create/drop databases/tables)")
			conn.Write([]byte("{\"status\":\"error\",\"message\":\"only master can create/drop databases/tables\"}\n"))
			continue
		}

		fmt.Printf("[%s] %s\n", req.FromSlave, req.Query)
		logger.WriteString(fmt.Sprintf("[%s] %s\n", req.FromSlave, req.Query))

		resp := shared.Response{}
		if req.IsSelect {
			rows, err := db.QueryRows(req.Query)
			if err != nil {
				resp.Status = "error"
				resp.Message = err.Error()
			} else {
				cols, _ := rows.Columns()
				resp.Header = cols
				for rows.Next() {
					colsVals := make([]interface{}, len(cols))
					colsPtrs := make([]interface{}, len(cols))
					for i := range colsVals {
						colsPtrs[i] = &colsVals[i]
					}
					rows.Scan(colsPtrs...)

					strRow := make([]interface{}, len(cols))
					for i, val := range colsVals {
						if b, ok := val.([]byte); ok {
							strRow[i] = string(b)
						} else {
							strRow[i] = val
						}
					}
					resp.Rows = append(resp.Rows, strRow)
				}
				resp.Status = "ok"
				resp.Message = "Select executed"
			}
		} else {
			affected, err := db.ExecQuery(req.Query)
			if err != nil {
				resp.Status = "error"
				resp.Message = err.Error()
			} else {
				resp.Status = "ok"
				resp.Message = fmt.Sprintf("Query executed successfully. Rows affected: %d", affected)
			}
		}
		respData, _ := json.Marshal(resp)
		conn.Write(append(respData, '\n'))
	}
}

func StartServer(db *shared.DBHandler) {
	logger, _ := os.OpenFile("master_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logger.Close()
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	fmt.Println("Master is listening on port 9000...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go handleSlave(conn, db, logger)
	}
}

func handleLocalQuery(req shared.Request, db *shared.DBHandler) shared.Response {
	resp := shared.Response{}

	if (strings.HasPrefix(strings.ToUpper(req.Query), "CREATE") || strings.HasPrefix(strings.ToUpper(req.Query), "DROP")) &&
		req.FromSlave != "master" {
		resp.Status = "error"
		resp.Message = "Only master can create/drop databases/tables"
		return resp
	}

	if req.IsSelect {
		rows, err := db.QueryRows(req.Query)
		if err != nil {
			resp.Status = "error"
			resp.Message = err.Error()
		} else {
			cols, _ := rows.Columns()
			resp.Header = cols
			for rows.Next() {
				colsVals := make([]interface{}, len(cols))
				colsPtrs := make([]interface{}, len(cols))
				for i := range colsVals {
					colsPtrs[i] = &colsVals[i]
				}
				rows.Scan(colsPtrs...)

				strRow := make([]interface{}, len(cols))
				for i, val := range colsVals {
					if b, ok := val.([]byte); ok {
						strRow[i] = string(b)
					} else {
						strRow[i] = val
					}
				}
				resp.Rows = append(resp.Rows, strRow)
			}
			resp.Status = "ok"
			resp.Message = "Select executed"
		}
	} else {
		affected, err := db.ExecQuery(req.Query)
		if err != nil {
			resp.Status = "error"
			resp.Message = err.Error()
		} else {
			resp.Status = "ok"
			resp.Message = fmt.Sprintf("Query executed successfully. Rows affected: %d", affected)
		}
	}
	return resp
}

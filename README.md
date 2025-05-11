# Database_Connection
# Distributed Database Project (Master-Slave Architecture)

This is a distributed SQL database system written in Go. It follows a **master-slave** architecture:
- The **master** node handles critical queries like `CREATE DATABASE`, `DROP TABLE`, etc.
- The **slaves** can connect to the master and execute standard queries (`SELECT`, `INSERT`, `UPDATE`, etc.).
- All communication is done over TCP using a custom protocol with JSON-formatted requests.

---

## 🏗️ Project Structure

```
distributed-db/
├── master/
│   ├── main.go        # Starts the master node and waits for queries
│   └── server.go      # TCP server handling requests from slaves
├── slave/
│   └── main.go        # Connects to the master and sends queries
├── shared/
│   └── db.go          # Database logic shared between master and slave
├── web/               # Simple GUI for the master
│   ├── index.html
│   ├── style.css
│   └── script.js
└── go.mod
```

---

## 🚀 How It Works

- Master listens on port `:9000` and waits for TCP connections from slave nodes.
- Queries are sent as JSON-encoded messages with a token for authorization.
- The master logs all queries from slaves in `master_log.txt`.
- `SELECT` query results are parsed and sent back with headers and rows.

---

## 🖥️ GUI (Web-based Interface)

The project includes a **simple web GUI** that can be used to send SQL queries to the master server.

```
web/
├── index.html   # Main interface
├── style.css    # UI styling
└── script.js    # Sends queries and handles results
```

You can serve this using any static file server or extend the master Go code to serve it.

---

## 🧪 Example Commands

From the master terminal:
```
> create database testdb;
> use testdb;
> create table users (id INT, name VARCHAR(50));
> insert into users values (1, 'ahmed');
> select * from users;
```

---

## ⚠️ Notes

- The master no longer assumes a default database. You must start by issuing `CREATE DATABASE` or `USE databasename`.
- Only the master is allowed to run DDL statements (CREATE, DROP).
- Slave queries are filtered and validated by the master.
- SELECT query results are returned in clean format (with proper UTF-8 decoding).

---

## 🛡️ Security

- Basic token-based validation is in place (`Token: "secret-token"`).
- All connections are handled over TCP. You may extend this with TLS or SSH tunnels.

---

## 📦 Requirements

- Go 1.18+
- MySQL server running (default on `127.0.0.1:3306` with user `root`/`rootroot`)
- Modern browser for GUI

---

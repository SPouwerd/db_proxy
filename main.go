package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	proxy_auth "github.com/OPEN-ICT-intergrator/database_proxy/auth"
	proxy_tpc "github.com/OPEN-ICT-intergrator/database_proxy/tpc"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	ln, err := net.Listen("tcp", ":5432")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	log.Println("Listening to TCP requests on port", ln.Addr().String())

	for {
		clientConn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(clientConn)
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	log.Println("New connection:", clientConn.RemoteAddr().String())

	db, err := connectToDatabase()
	if db == nil {
		log.Println("Error connecting to database:", err)
		return
	}

	username, password, projectName, err := proxy_tpc.ParseConnectionString(clientConn)
	if err != nil {
		log.Println("Error parsing connectionString:", err)
		return
	}

	exists, err := proxy_auth.AuthenticateUser(username, password, db)
	if err != nil || !exists {
		log.Println("User does not exist")
		return
	}

	project, err := proxy_auth.FindProjectByShortName(projectName, db)
	if err != nil {
		log.Println("Project does not exist: ", projectName)
		return
	}

	if !proxy_auth.IsIPWhitelisted(clientConn.RemoteAddr().String(), project) {
		log.Println("Client's IP is not whitelisted for this project")
		return
	}

	serverConn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println("Error connecting to PostgreSQL server:", err)
		return
	}
	defer serverConn.Close()

	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}

func getConnectionString() string {
	var _ = godotenv.Load(".env")
	return fmt.Sprintf(
		"host=localhost port=8080 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB_NAME"))
}
func connectToDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}

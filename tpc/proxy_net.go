package proxy_net

import (
	"log"
	"net"
	"strings"
)

func readBuffer(conn net.Conn) (string, error) {
	buffer := make([]byte, 8*12) // TODO: test with bigger connection strings and set max length on username & password
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

func ParseConnectionString(clientConn net.Conn) (username, password, projectName string, err error) {

	connStr, err := readBuffer(clientConn)
	if err != nil {
		log.Println("Error reading connection string:", err)
	}
	log.Print("Received connection string: \n", connStr, "\n")

	split := strings.Split(connStr, "database")
	// HOURS WASTED 3 - to split the user and database string in a readable way
	username = split[0][13:]
	projectName = split[1][:8]

	nextBuff, err := readBuffer(clientConn)
	if err != nil {
		log.Println("Error reading next BUffer:", err)
	}
	log.Default().Print("Received next BUffer: \n", nextBuff, "\n")

	return username, password, projectName, nil
}

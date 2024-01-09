package proxy_auth

import (
	"database/sql"
	"log"
)

type Project struct {
	// id            string
	// supervisor_id string
	// name          string
	// short_name    string
	// description   string
	whitelist_ip []string
	// created_at    string
}

func AuthenticateUser(username, password string, db *sql.DB) (bool, error) {
	log.Println("Authenticating user:", username, "with password: ", password)
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1 AND password=$2)", username, password).Scan(&exists)
	return exists, err
}

func FindProjectByShortName(shortName string, db *sql.DB) (project Project, err error) {
	err = db.QueryRow("SELECT * FROM projects WHERE short_name=$1", shortName).Scan(&project)
	return project, err
}

func IsIPWhitelisted(clientIP string, project Project) bool {
	for _, ip := range project.whitelist_ip {
		if clientIP == ip {
			return true
		}
	}
	return false
}

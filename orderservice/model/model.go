package model

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Database *sql.DB
}


func (s *Server) CreateOrder(guid string, timestamp int, cost int) {
	query := "INSERT INTO `order` (id, created_timestamp, cost) VALUES (?, ?, ?)"
	_, err := s.Database.Exec(query, guid, 100000, 1)
	if err != nil {
		log.WithField("create_order", "failed")
	}
}

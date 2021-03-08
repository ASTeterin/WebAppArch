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
	_, err := s.Database.Exec(query, guid, timestamp, cost)
	if err != nil {
		log.WithField("create_order", "failed")
	}
}

func (s *Server) DeleteOrder(id string) error {
	query := "DELETE FROM `order` WHERE id = ?"
	_, err := s.Database.Exec(query, id)
	return err
}

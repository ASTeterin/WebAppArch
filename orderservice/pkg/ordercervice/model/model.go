package model

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Database *sql.DB
}

type orderResponse struct {
	order
	OrderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type order struct {
	Id    string `json:"id"`
	//menuItems []menuItem
}


func (s *Server) CreateOrder(guid string, timestamp int, cost int) {
	query := "INSERT INTO `order` (id, created_timestamp, cost) VALUES (?, ?, ?)"
	_, err := s.Database.Exec(query, guid, timestamp, cost)
	if err != nil {
		log.WithField("create_order", "failed")
	}
}

func (s *Server) DeleteOrder(id string) {
	query := "DELETE FROM `order` WHERE id = ?"
	_, err := s.Database.Exec(query, id)
	if err != nil {
		log.WithField("create_order", "failed")
	}
}

func (s *Server) GetOrders() []orderResponse {
	query := "SELECT * FROM `order`"
	rows, err := s.Database.Query(query)
	if err != nil {
		log.WithField("select_order", "failed")
	}
	defer rows.Close()
	var orders []orderResponse
	var order orderResponse

	for rows.Next() {
		err := rows.Scan(&order.Id, &order.OrderedAtTimeStamp, &order.Cost)
		if err != nil{
			fmt.Println(err)
			continue
		}
		orders = append(orders, order)
	}
	return orders
}

func (s *Server) updateOrder() {

}

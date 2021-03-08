package model

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Database *sql.DB
}

type OrderResponse struct {
	Order
	OrderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type Order struct {
	Id    string `json:"id"`
	//MenuItems []MenuItem
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

func (s *Server) GetOrders() []OrderResponse{
	query := "SELECT * FROM `order`"
	rows, err := s.Database.Query(query)
	if err != nil {
		log.WithField("select_order", "failed")
	}
	defer rows.Close()
	var orders []OrderResponse
	var order OrderResponse
	/*var id string
	var timestamp int
	var cost int
*/
	for rows.Next() {
		err := rows.Scan(&order.Id, &order.OrderedAtTimeStamp, &order.Cost)
		if err != nil{
			fmt.Println(err)
			continue
		}
		orders = append(orders, order)
		//fmt.Println(id,  timestamp, cost)
	}
	return orders
}

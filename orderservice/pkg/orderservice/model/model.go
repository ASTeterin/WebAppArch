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
	menuItems []menuItem
}

type menuItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type MenuItems struct {
	MenuItems  []menuItem `json:"menuItems"`
}


func (s *Server) CreateOrder(guid string, timestamp int, cost int, menuItems []menuItem) {
	query := "INSERT INTO `order` (id, created_timestamp, cost) VALUES (?, ?, ?)"
	_, err := s.Database.Exec(query, guid, timestamp, cost)
	if err != nil {
		log.WithField("create_order", "failed")
	}
	for _, item := range  menuItems{

		query = "INSERT INTO `menu_item` (idmenu_item, quantity) VALUES (?, ?)"
		result, err := s.Database.Exec(query, item.Id, item.Quantity)
		menuItemId, _ := result.LastInsertId()
		if err != nil {
			log.WithField("create_order", "failed")
		}
		fmt.Println(item)
		query = "INSERT INTO `item_in_order` (order_id, menu_item_id) VALUES (?, ?)"
		_, err = s.Database.Exec(query, guid, int(menuItemId))
		if err != nil {
			log.WithField("create_order", "failed")
		}
	}
}

func (s *Server) DeleteOrder(id string) {
	query := "DELETE FROM `order` WHERE id = ?"
	_, err := s.Database.Exec(query, id)
	if err != nil {
		log.WithField("delete_order", "failed")
	}

	query = "SELECT menu_item_id FROM `item_in_order` WHERE order_id = ?"
	rows, err := s.Database.Query(query, id)
	if err != nil {
		fmt.Println("no rows")
	}

	defer rows.Close()

	for rows.Next(){
		var menuItemId int
 		err = rows.Scan(&menuItemId)
		fmt.Print(menuItemId, ' ')
 		if err != nil {
			log.WithField("delete menu_item", "failed")
 			return
		}
		query = "DELETE FROM `menu_item` WHERE id = ?"
		_, err = s.Database.Exec(query, menuItemId)
		if err != nil {
			log.WithField("delete menu_item", "failed")
			return
		}

	}

	query = "DELETE FROM `item_in_order` WHERE order_id = ?"
	_, err = s.Database.Exec(query, id)
	if err != nil {
		log.WithField("delete item_in_order", "failed")

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

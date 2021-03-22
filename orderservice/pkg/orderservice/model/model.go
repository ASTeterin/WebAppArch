package model

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Database *sql.DB
}

type orderDTO struct {
	Id    string `json:"id"`
	OrderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type OrderResponse struct {
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

func deleteOrderParam(s *Server, id string) {
	query := "SELECT menu_item_id FROM `item_in_order` WHERE order_id = ?"
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

func (s *Server) DeleteOrder(id string) {
	query := "DELETE FROM `order` WHERE id = ?"
	_, err := s.Database.Exec(query, id)
	if err != nil {
		log.WithField("delete_order", "failed")
	}
	deleteOrderParam(s, id)
}

func (s *Server) GetOrders() []OrderResponse {
	query := "SELECT * FROM `order`"
	rows, err := s.Database.Query(query)
	if err != nil {
		log.WithField("select_order", "failed")
	}
	defer rows.Close()
	var orders []OrderResponse
	var order OrderResponse

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

func (s *Server) GetOrder(id string)  OrderResponse{
	//var orderResponse orderResponse
	/*var time string
	var cost int
	query := "SELECT creared_timestamp, cost FROM `order` WHERE order_id = ?"
	err := s.Database.QueryRow(query, id).Scan(&time, &cost)
	if err != nil {
		log.WithField("find_order", "failed")
		return
	}
	fmt.Printf(string(cost), ' ')
*/
	var orderResponse OrderResponse
	query := "SELECT menu_item_id FROM item_in_order WHERE order_id = ?"
	rows, err := s.Database.Query(query, id)
	if err != nil {
		log.WithField("select_order", "failed")
	}
	defer rows.Close()

	menuItems := make([]menuItem, 0)
	for rows.Next(){
		//fmt.Println("11111")
		var currentMenuItem menuItem
		var menuItemId int
		err = rows.Scan(&menuItemId)
		fmt.Print(menuItemId, ' ')
		if err != nil {
			log.WithField("find menu_item", "failed")
			return orderResponse
		}
		query = "SELECT idmenu_item, quantity FROM `menu_item` WHERE id = ?"
		err = s.Database.QueryRow(query, menuItemId).Scan(&currentMenuItem.Id, &currentMenuItem.Quantity)
		if err != nil {
			log.WithField("delete menu_item", "failed")
			return orderResponse
		}
		menuItems = append(menuItems, currentMenuItem)
		fmt.Println(currentMenuItem)
	}
	orderResponse.menuItems = menuItems
	orderResponse.Id = id
	return orderResponse
}

func createMenuItemsInOrder(menuItems []menuItem, id string, s *Server) {
	for _, item := range  menuItems{

		query := "INSERT INTO `menu_item` (idmenu_item, quantity) VALUES (?, ?)"
		result, err := s.Database.Exec(query, item.Id, item.Quantity)
		menuItemId, _ := result.LastInsertId()
		if err != nil {
			log.WithField("create_order", "failed")
		}
		fmt.Println(item)
		query = "INSERT INTO `item_in_order` (order_id, menu_item_id) VALUES (?, ?)"
		_, err = s.Database.Exec(query, id, int(menuItemId))
		if err != nil {
			log.WithField("create_order", "failed")
		}
	}
}

func (s *Server) UpdateOrder(id string, timestamp int, cost int,  menuItems []menuItem) {
	fmt.Printf(id)
	query := "UPDATE `order` SET created_timestamp = ?, cost = ? WHERE id = ?"
	_, err := s.Database.Exec(query, timestamp, cost, id)
	if err != nil {
		log.WithField("create_order", "failed")
	}

	deleteOrderParam(s, id)
	createMenuItemsInOrder(menuItems, id, s)
}



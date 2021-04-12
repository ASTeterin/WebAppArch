package model

import (
	"database/sql"
	"fmt"
)

type DBServer struct {
	Database *sql.DB
}

type OrderResponse struct {
	order
	OrderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type order struct {
	Id        string `json:"id"`
	menuItems []menuItem
}

type menuItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type MenuItems struct {
	MenuItems []menuItem `json:"menuItems"`
}

type OrderServiceInterface interface {
	CreateOrder(guid string, timestamp int, cost int, menuItems []menuItem) error
	DeleteOrder(id string) error
	GetOrders() ([]OrderResponse, error)
	GetOrder(id string) (OrderResponse, error)
	UpdateOrder(id string, timestamp int, cost int, menuItems []menuItem) error
}

func CreateOrderServiceInterface(db *sql.DB) OrderServiceInterface {
	return &DBServer{Database: db}
}

func (s *DBServer) CreateOrder(guid string, timestamp int, cost int, menuItems []menuItem) error {
	query := "INSERT INTO `order` (id, created_timestamp, cost) VALUES (?, ?, ?)"
	_, err := s.Database.Exec(query, guid, timestamp, cost)
	if err != nil {
		return err
	}
	for _, item := range menuItems {

		query = "INSERT INTO `menu_item` (idmenu_item, quantity) VALUES (?, ?)"
		result, err := s.Database.Exec(query, item.Id, item.Quantity)
		menuItemId, _ := result.LastInsertId()
		if err != nil {
			return err
		}
		fmt.Println(item)
		query = "INSERT INTO `item_in_order` (order_id, menu_item_id) VALUES (?, ?)"
		_, err = s.Database.Exec(query, guid, int(menuItemId))
		if err != nil {
			return err
		}
	}
	return err
}

func deleteOrderParam(s *DBServer, id string) error {
	query := "SELECT menu_item_id FROM `item_in_order` WHERE order_id = ?"
	rows, err := s.Database.Query(query, id)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var menuItemId int
		err = rows.Scan(&menuItemId)
		fmt.Print(menuItemId, ' ')
		if err != nil {
			return err
		}
		query = "DELETE FROM `menu_item` WHERE id = ?"
		_, err = s.Database.Exec(query, menuItemId)
		if err != nil {
			return err
		}
	}

	query = "DELETE FROM `item_in_order` WHERE order_id = ?"
	_, err = s.Database.Exec(query, id)
	return err
}

func (s *DBServer) DeleteOrder(id string) error {
	query := "DELETE FROM `order` WHERE id = ?"
	_, err := s.Database.Exec(query, id)
	if err != nil {
		return err
	}
	err = deleteOrderParam(s, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBServer) GetOrders() ([]OrderResponse, error) {
	query := "SELECT * FROM `order`"
	rows, err := s.Database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []OrderResponse
	var order OrderResponse

	for rows.Next() {
		err := rows.Scan(&order.Id, &order.OrderedAtTimeStamp, &order.Cost)
		if err != nil {
			return nil, err
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *DBServer) GetOrder(id string) (OrderResponse, error) {
	var orderResponse OrderResponse
	query := "SELECT menu_item_id FROM item_in_order WHERE order_id = ?"
	rows, err := s.Database.Query(query, id)
	if err != nil {
		return OrderResponse{}, err
	}
	defer rows.Close()

	menuItems := make([]menuItem, 0)
	for rows.Next() {
		var currentMenuItem menuItem
		var menuItemId int
		err = rows.Scan(&menuItemId)
		if err != nil {
			return OrderResponse{}, err
		}
		query = "SELECT idmenu_item, quantity FROM `menu_item` WHERE id = ?"
		err = s.Database.QueryRow(query, menuItemId).Scan(&currentMenuItem.Id, &currentMenuItem.Quantity)
		if err != nil {
			return OrderResponse{}, nil
		}
		menuItems = append(menuItems, currentMenuItem)
	}
	orderResponse.menuItems = menuItems
	orderResponse.Id = id
	return orderResponse, nil
}

func createMenuItemsInOrder(menuItems []menuItem, id string, s *DBServer) error {
	for _, item := range menuItems {

		query := "INSERT INTO `menu_item` (idmenu_item, quantity) VALUES (?, ?)"
		result, err := s.Database.Exec(query, item.Id, item.Quantity)
		menuItemId, _ := result.LastInsertId()
		if err != nil {
			return err
		}
		fmt.Println(item)
		query = "INSERT INTO `item_in_order` (order_id, menu_item_id) VALUES (?, ?)"
		_, err = s.Database.Exec(query, id, int(menuItemId))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DBServer) UpdateOrder(id string, timestamp int, cost int, menuItems []menuItem) error {
	fmt.Printf(id)
	query := "UPDATE `order` SET created_timestamp = ?, cost = ? WHERE id = ?"
	_, err := s.Database.Exec(query, timestamp, cost, id)
	if err != nil {
		return err
	}

	err = deleteOrderParam(s, id)
	if err != nil {
		return err
	}
	err = createMenuItemsInOrder(menuItems, id, s)
	if err != nil {
		return err
	}
	return nil
}

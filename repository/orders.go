package repository

import (
	"context"

	"order_system/model"

	"github.com/jackc/pgx/v5"
)

type OrderRepository interface {
	Ping(ctx context.Context) error
	Store(ctx context.Context, order model.Order) error
	GetOrderStatus(ctx context.Context, orderID int) (string, error)
	UpdateOrderStatus(ctx context.Context, orderID int, status string) error
}

type Orders struct {
	MasterDB *pgx.Conn
	SlaveDB  *pgx.Conn
}

func NewOrderRepository(masterDB, slaveDB *pgx.Conn) OrderRepository {
	return &Orders{
		MasterDB: masterDB,
		SlaveDB:  slaveDB,
	}
}

func (o *Orders) Ping(ctx context.Context) error {
	err := o.MasterDB.Ping(ctx)
	if err != nil {
		return err
	}
	err = o.SlaveDB.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (o *Orders) Store(ctx context.Context, order model.Order) error {
	_, err := o.MasterDB.Exec(context.Background(), "INSERT INTO orders (user_id, product, count, address, status) VALUES ($1, $2, $3, $4, $5)", order.UserID, order.Product, order.Count, order.Address, model.OrderStatusPending)
	if err != nil {
		return err
	}
	return nil
}

func (o *Orders) GetOrderStatus(ctx context.Context, orderID int) (string, error) {
	var status string
	err := o.SlaveDB.QueryRow(context.Background(), "SELECT status FROM orders WHERE id = $1", orderID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (o *Orders) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	_, err := o.MasterDB.Exec(context.Background(), "UPDATE orders SET status = $1 WHERE id = $2", status, orderID)
	if err != nil {
		return err
	}
	return nil
}

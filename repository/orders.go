package repository

import (
	"context"

	"order_system/model"

	"github.com/jackc/pgx/v5"
)

type OrderRepository interface {
	Ping(ctx context.Context) error
	Store(ctx context.Context, order model.Order) (int, error)
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

func (o *Orders) Store(ctx context.Context, order model.Order) (int, error) {
	var id int
	err := o.MasterDB.QueryRow(ctx, "INSERT INTO orders (user_id, product, count, address) VALUES ($1, $2, $3, $4) RETURNING id", order.UserID, order.Product, order.Count, order.Address).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (o *Orders) GetOrderStatus(ctx context.Context, orderID int) (string, error) {
	var status string
	err := o.SlaveDB.QueryRow(ctx, "SELECT status FROM orders WHERE id = $1", orderID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (o *Orders) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	_, err := o.MasterDB.Exec(ctx, "UPDATE orders SET status = $1 WHERE id = $2", status, orderID)
	if err != nil {
		return err
	}
	return nil
}

package http

import (
	"context"
	"fmt"
	"log"

	"order_system/internal/config"
	"order_system/repository"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e               *echo.Echo
	OrderRepository repository.OrderRepository
	Config          *config.Config
}

func NewServer() *Server {
	cfg, err := config.LoadConfig("/app/config.yaml")
	if err != nil {
		log.Fatalf("Unable to load config: %v\n", err)
	}
	return &Server{
		e:      echo.New(),
		Config: cfg,
	}
}

func (s *Server) Start(port string) error {
	fmt.Println("Starting server on port", s.Config)
	s.withPostgres()
	v1 := s.e.Group("/v1")
	v1.POST("/orders", s.CreateOrder)
	v1.GET("/orders/:id/status", s.GetOrderStatus)
	v1.PUT("/orders/:id", s.UpdateOrderStatus)
	v1.GET("/startup", s.Startup)
	v1.GET("/readiness", s.Readiness)
	v1.GET("/liveness", s.Liveness)
	return s.e.Start(":" + port)
}

func (s *Server) withPostgres() error {
	masterDB, err := pgx.Connect(context.Background(),
		"host="+s.Config.Postgres.MasterDB.Host+
			" port="+s.Config.Postgres.MasterDB.Port+
			" user="+s.Config.Postgres.MasterDB.User+
			" password="+s.Config.Postgres.MasterDB.Password+
			" dbname="+s.Config.Postgres.MasterDB.DBName)
	if err != nil {
		log.Fatalf("Unable to connect to master database: %v\n", err)
	}
	slaveDB, err := pgx.Connect(context.Background(),
		"host="+s.Config.Postgres.SlaveDB.Host+
			" port="+s.Config.Postgres.SlaveDB.Port+
			" user="+s.Config.Postgres.SlaveDB.User+
			" password="+s.Config.Postgres.SlaveDB.Password+
			" dbname="+s.Config.Postgres.SlaveDB.DBName)
	if err != nil {
		log.Fatalf("Unable to connect to slave database: %v\n", err)
	}

	s.OrderRepository = repository.NewOrderRepository(masterDB, slaveDB)
	return nil
}

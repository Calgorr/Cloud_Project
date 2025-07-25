package http

import (
	"net/http"
	"strconv"

	"order_system/internal/config"
	"order_system/model"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreateOrder(c echo.Context) error {
	var order model.Order
	err := c.Bind(&order)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	id, err := s.OrderRepository.Store(c.Request().Context(), order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusCreated, "Order created"+strconv.Itoa(id))
}

func (s *Server) GetOrderStatus(c echo.Context) error {
	orderID := c.Param("id")
	orderIDs, _ := strconv.Atoi(orderID)
	status, err := s.OrderRepository.GetOrderStatus(c.Request().Context(), orderIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, status)

}

func (s *Server) UpdateOrderStatus(c echo.Context) error {
	OrderID := struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}{}
	err := c.Bind(&OrderID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	err = s.OrderRepository.UpdateOrderStatus(c.Request().Context(), OrderID.ID, OrderID.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, "Order status updated")
}

func (s *Server) Startup(c echo.Context) error {
	cfg, err := config.LoadConfig("/app/config.yaml")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to load configurations")
	}
	s.Config = cfg
	return c.JSON(http.StatusOK, "Service started successfully")
}

func (s *Server) Liveness(c echo.Context) error {
	return c.JSON(http.StatusOK, "Service alive")
}

func (s *Server) Readiness(c echo.Context) error {
	err := s.OrderRepository.Ping(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, "Service not ready")
	}
	return c.JSON(http.StatusOK, "Service ready")
}

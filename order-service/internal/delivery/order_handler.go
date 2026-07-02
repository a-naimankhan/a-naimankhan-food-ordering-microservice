package delivery

import (
	"net/http"
	"order-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderHandler struct {
	service  domain.OrderService
	validate *validator.Validate
}

func NewOrderHandler(svc domain.OrderService) *OrderHandler {
	return &OrderHandler{service: svc, validate: validator.New()}
}

type OrderRequest struct {
	CustomerID string  `json:"customer_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Status     string  `json:"status" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		ErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), req.CustomerID, req.Amount, req.Status)
	if err != nil {
		ErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponce(c, http.StatusCreated, order)
}

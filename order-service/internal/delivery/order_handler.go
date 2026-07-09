package delivery

import (
	"errors"
	"net/http"
	"order-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cid, err := uuid.Parse(req.CustomerID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	orderModel := &domain.Order{
		CustomerID: cid,
		Amount:     req.Amount,
		Status:     req.Status,
	}

	created, err := h.service.CreateOrder(c.Request.Context(), orderModel)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCustomerID) || errors.Is(err, domain.ErrInvalidAmount) || errors.Is(err, domain.ErrInvalidStatus) {
			ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, created)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		ErrorResponse(c, http.StatusNotFound, "id is required")
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "not valid format id")
		return
	}
	order, err := h.service.GetOrder(c.Request.Context(), parsedID)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "order not found")
		return
	}
	SuccessResponse(c, http.StatusOK, order)
}

func (h *OrderHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

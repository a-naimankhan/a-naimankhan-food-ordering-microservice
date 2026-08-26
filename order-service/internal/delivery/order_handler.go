package delivery

import (
	"errors"
	"net/http"
	"order-service/internal/domain"
	"order-service/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type OrderHandler struct {
	service  domain.OrderService
	validate *validator.Validate
	log      *logger.Logger
}

func NewOrderHandler(svc domain.OrderService) *OrderHandler {
	return &OrderHandler{
		service:  svc,
		validate: validator.New(),
		log:      logger.GetLogger(),
	}
}

type OrderRequest struct {
	CustomerID string  `json:"customer_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Status     string  `json:"status" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	h.log.Debug("HTTP: POST /api/v1/orders received | RemoteAddr=%s", c.RemoteIP())
	
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("HTTP: Invalid JSON request | Error: %v", err)
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.log.Debug("HTTP: Request parsed | CustomerID=%s | Amount=%.2f | Status=%s", 
		req.CustomerID, req.Amount, req.Status)

	if err := h.validate.Struct(req); err != nil {
		h.log.Error("HTTP: Validation failed | Error: %v", err)
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cid, err := uuid.Parse(req.CustomerID)
	if err != nil {
		h.log.Error("HTTP: Invalid customer ID format | CustomerID=%s | Error: %v", req.CustomerID, err)
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	orderModel := &domain.Order{
		CustomerID: cid,
		Amount:     req.Amount,
		Status:     req.Status,
	}

	h.log.Info("HTTP: Creating order | CustomerID=%s | Amount=%.2f", cid, req.Amount)
	created, err := h.service.CreateOrder(c.Request.Context(), orderModel)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCustomerID) || errors.Is(err, domain.ErrInvalidAmount) || errors.Is(err, domain.ErrInvalidStatus) {
			h.log.Error("HTTP: Validation error from service | Error: %v", err)
			ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		h.log.Error("HTTP: Internal server error | Error: %v", err)
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Info("✅ HTTP: Order created successfully | ID=%s | StatusCode=201", created.ID)
	SuccessResponse(c, http.StatusCreated, created)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	h.log.Debug("HTTP: GET /api/v1/orders/:id received | ID=%s | RemoteAddr=%s", id, c.RemoteIP())

	if id == "" {
		h.log.Error("HTTP: Order ID is required")
		ErrorResponse(c, http.StatusNotFound, "id is required")
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.log.Error("HTTP: Invalid order ID format | ID=%s | Error: %v", id, err)
		ErrorResponse(c, http.StatusBadRequest, "not valid format id")
		return
	}
	
	h.log.Debug("HTTP: Fetching order from service | ID=%s", parsedID)
	order, err := h.service.GetOrder(c.Request.Context(), parsedID)
	if err != nil {
		h.log.Error("HTTP: Order not found | ID=%s | Error: %v", parsedID, err)
		ErrorResponse(c, http.StatusNotFound, "order not found")
		return
	}
	
	h.log.Info("✅ HTTP: Order retrieved successfully | ID=%s | StatusCode=200", parsedID)
	SuccessResponse(c, http.StatusOK, order)
}

func (h *OrderHandler) Ping(c *gin.Context) {
	h.log.Debug("HTTP: GET /api/v1/ping received | RemoteAddr=%s", c.RemoteIP())
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	h.log.Debug("HTTP: Ping response sent | StatusCode=200")
}

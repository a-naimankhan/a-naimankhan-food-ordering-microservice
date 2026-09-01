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
	h.log.Debug("POST /api/v1/orders received", "remote_addr", c.RemoteIP())

	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid JSON request", "error", err.Error())
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.log.Debug("Request parsed", "customer_id", req.CustomerID, "amount", req.Amount, "status", req.Status)

	//if err := h.validate.Struct(req); err != nil {
	//	h.log.Error("Validation failed", "error", err.Error())
	//	ErrorResponse(c, http.StatusBadRequest, err.Error())
	//	return
	//}

	cid, err := uuid.Parse(req.CustomerID)
	if err != nil {
		h.log.Error("Invalid customer ID format", "customer_id", req.CustomerID, "error", err.Error())
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	orderModel := &domain.Order{
		CustomerID: cid,
		Amount:     req.Amount,
		Status:     req.Status,
	}

	h.log.Info("Creating order", "customer_id", cid, "amount", req.Amount)
	created, err := h.service.CreateOrder(c.Request.Context(), orderModel)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCustomerID) || errors.Is(err, domain.ErrInvalidAmount) || errors.Is(err, domain.ErrInvalidStatus) {
			h.log.Error("Validation error from service", "error", err.Error())
			ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		h.log.Error("Internal server error", "error", err.Error())
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Info("Order created successfully", "id", created.ID, "status_code", 201)
	SuccessResponse(c, http.StatusCreated, created)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	h.log.Debug("GET /api/v1/orders/:id received", "id", id, "remote_addr", c.RemoteIP())

	if id == "" {
		h.log.Error("Order ID is required")
		ErrorResponse(c, http.StatusNotFound, "id is required")
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.log.Error("Invalid order ID format", "id", id, "error", err.Error())
		ErrorResponse(c, http.StatusBadRequest, "not valid format id")
		return
	}

	h.log.Debug("Fetching order from service", "id", parsedID)
	order, err := h.service.GetOrder(c.Request.Context(), parsedID)
	if err != nil {
		h.log.Error("Order not found", "id", parsedID, "error", err.Error())
		ErrorResponse(c, http.StatusNotFound, "order not found")
		return
	}

	h.log.Info("Order retrieved successfully", "id", parsedID, "status_code", 200)
	SuccessResponse(c, http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	h.log.Debug("PATCH /api/v1/orders/:id/status received", "id", id, "remote_addr", c.RemoteIP())

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.log.Error("Invalid order ID format", "id", id, "error", err.Error())
		ErrorResponse(c, http.StatusBadRequest, "not valid format id")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid JSON request for status update", "error", err.Error())
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.UpdateOrderStatus(c.Request.Context(), parsedID, req.Status)
	if err != nil {
		h.log.Error("Failed to update order status", "id", parsedID, "error", err.Error())
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Info("Order status updated successfully", "id", parsedID, "status", req.Status)
	SuccessResponse(c, http.StatusOK, gin.H{"status": req.Status})
	return

}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	h.log.Debug("PATCH /api/v1/orders/:id/cancel received", "id", id, "remote_addr", c.RemoteIP())

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.log.Error("Invalid order ID format", "id", id, "error", err.Error())
		ErrorResponse(c, http.StatusBadRequest, "not valid format id")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid JSON request", "error", err.Error())
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.CancelOrder(c.Request.Context(), parsedID, req.Reason)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			h.log.Error("Order not found", "id", parsedID, "error", err.Error())
			ErrorResponse(c, http.StatusNotFound, "order not found")
			return
		}

		if errors.Is(err, domain.ErrCannotCancelOrder) {
			h.log.Error("Order cannot be cancelled", "id", parsedID, "error", err.Error())
			ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}

		h.log.Error("Internal error", "id", parsedID, "error", err.Error())
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Info("Order cancelled successfully", "id", parsedID)
	SuccessResponse(c, http.StatusOK, gin.H{"status": domain.StatusCancelled})
}

func (h *OrderHandler) Ping(c *gin.Context) {
	h.log.Debug("GET /api/v1/ping received", "remote_addr", c.RemoteIP())
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	h.log.Debug("Ping response sent", "status_code", 200)
}

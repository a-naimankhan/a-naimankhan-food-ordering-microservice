package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"order-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeService struct {
	CreateOrderFn       func(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetOrderFn          func(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	UpdateOrderStatusFn func(ctx context.Context, id uuid.UUID, status string) error
	CancelOrderFn       func(ctx context.Context, id uuid.UUID, reason string) error
}

func (f *fakeService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	if f.CreateOrderFn == nil {
		return nil, nil
	}
	return f.CreateOrderFn(ctx, order)
}

func (f *fakeService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	if f.GetOrderFn == nil {
		return nil, nil
	}
	return f.GetOrderFn(ctx, id)
}

func (f *fakeService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	if f.UpdateOrderStatusFn == nil {
		return nil
	}
	return f.UpdateOrderStatusFn(ctx, id, status)
}

func (f *fakeService) CancelOrder(ctx context.Context, id uuid.UUID, reason string) error {
	if f.CancelOrderFn == nil {
		return nil
	}
	return f.CancelOrderFn(ctx, id, reason)
}

func setupRouter(svc domain.OrderService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	h := NewOrderHandler(svc)

	r.POST("/orders", h.CreateOrder)
	r.GET("/orders/:id", h.GetOrder)
	r.PATCH("/orders/:id/status", h.UpdateOrderStatus)
	r.PATCH("/orders/:id/cancel", h.CancelOrder)

	return r
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	customerID := uuid.New()
	orderID := uuid.New()

	tests := []struct {
		name           string
		body           string
		service        *fakeService
		wantStatus     int
		wantMessage    string
		assertResponse func(t *testing.T, body []byte)
	}{
		{
			name: "success",
			body: jsonBody(t, map[string]interface{}{
				"customer_id": customerID.String(),
				"amount":      12.5,
				"status":      domain.StatusPending,
			}),
			service: &fakeService{
				CreateOrderFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
					assert.Equal(t, customerID, order.CustomerID)
					assert.Equal(t, 12.5, order.Amount)
					assert.Equal(t, domain.StatusPending, order.Status)

					return &domain.Order{
						ID:         orderID,
						CustomerID: order.CustomerID,
						Amount:     order.Amount,
						Status:     order.Status,
						CreatedAt:  time.Date(2026, 8, 31, 10, 0, 0, 0, time.UTC),
					}, nil
				},
			},
			wantStatus: http.StatusCreated,
			assertResponse: func(t *testing.T, body []byte) {
				var got domain.Order
				require.NoError(t, json.Unmarshal(body, &got))
				assert.Equal(t, orderID, got.ID)
				assert.Equal(t, customerID, got.CustomerID)
				assert.Equal(t, 12.5, got.Amount)
				assert.Equal(t, domain.StatusPending, got.Status)
			},
		},
		{
			name:        "invalid json",
			body:        `{"customer_id":`,
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "unexpected EOF",
		},
		{
			name: "missing required fields",
			body: jsonBody(t, map[string]interface{}{
				"customer_id": customerID.String(),
			}),
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Error:Field validation",
		},
		{
			name: "invalid customer id",
			body: jsonBody(t, map[string]interface{}{
				"customer_id": "not-a-uuid",
				"amount":      12.5,
				"status":      domain.StatusPending,
			}),
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "invalid UUID",
		},
		{
			name: "domain validation error",
			body: jsonBody(t, map[string]interface{}{
				"customer_id": customerID.String(),
				"amount":      -1,
				"status":      domain.StatusPending,
			}),
			service: &fakeService{
				CreateOrderFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
					return nil, domain.ErrInvalidAmount
				},
			},
			wantStatus:  http.StatusBadRequest,
			wantMessage: domain.ErrInvalidAmount.Error(),
		},
		{
			name: "unexpected service error",
			body: jsonBody(t, map[string]interface{}{
				"customer_id": customerID.String(),
				"amount":      12.5,
				"status":      domain.StatusPending,
			}),
			service: &fakeService{
				CreateOrderFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
					return nil, errors.New("service failure")
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "service failure",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.service)

			w := performRequest(router, http.MethodPost, "/orders", tt.body)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.assertResponse != nil {
				tt.assertResponse(t, w.Body.Bytes())
				return
			}
			assertErrorMessageContains(t, w.Body.Bytes(), tt.wantMessage)
		})
	}
}

func TestOrderHandler_GetOrder(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()

	tests := []struct {
		name           string
		idParam        string
		service        *fakeService
		wantStatus     int
		wantMessage    string
		assertResponse func(t *testing.T, body []byte)
	}{
		{
			name:    "success",
			idParam: orderID.String(),
			service: &fakeService{
				GetOrderFn: func(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
					assert.Equal(t, orderID, id)
					return &domain.Order{
						ID:         orderID,
						CustomerID: customerID,
						Amount:     42,
						Status:     domain.StatusDelivered,
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			assertResponse: func(t *testing.T, body []byte) {
				var got domain.Order
				require.NoError(t, json.Unmarshal(body, &got))
				assert.Equal(t, orderID, got.ID)
				assert.Equal(t, customerID, got.CustomerID)
				assert.Equal(t, 42.0, got.Amount)
				assert.Equal(t, domain.StatusDelivered, got.Status)
			},
		},
		{
			name:        "invalid uuid",
			idParam:     "not-a-uuid",
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "not valid format id",
		},
		{
			name:    "service returns not found",
			idParam: orderID.String(),
			service: &fakeService{
				GetOrderFn: func(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
					return nil, domain.ErrOrderNotFound
				},
			},
			wantStatus:  http.StatusNotFound,
			wantMessage: "order not found",
		},
		{
			name:    "service returns unexpected error",
			idParam: orderID.String(),
			service: &fakeService{
				GetOrderFn: func(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
					return nil, errors.New("database failure")
				},
			},
			wantStatus:  http.StatusNotFound,
			wantMessage: "order not found",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.service)

			w := performRequest(router, http.MethodGet, "/orders/"+tt.idParam, "")

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.assertResponse != nil {
				tt.assertResponse(t, w.Body.Bytes())
				return
			}
			assertErrorMessageContains(t, w.Body.Bytes(), tt.wantMessage)
		})
	}
}

func TestOrderHandler_UpdateOrderStatus(t *testing.T) {
	orderID := uuid.New()

	tests := []struct {
		name           string
		idParam        string
		body           string
		service        *fakeService
		wantStatus     int
		wantMessage    string
		assertResponse func(t *testing.T, body []byte)
	}{
		{
			name:    "success",
			idParam: orderID.String(),
			body: jsonBody(t, map[string]interface{}{
				"status": domain.StatusDelivered,
			}),
			service: &fakeService{
				UpdateOrderStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
					assert.Equal(t, orderID, id)
					assert.Equal(t, domain.StatusDelivered, status)
					return nil
				},
			},
			wantStatus: http.StatusOK,
			assertResponse: func(t *testing.T, body []byte) {
				var got map[string]string
				require.NoError(t, json.Unmarshal(body, &got))
				assert.Equal(t, domain.StatusDelivered, got["status"])
			},
		},
		{
			name:        "invalid uuid",
			idParam:     "not-a-uuid",
			body:        jsonBody(t, map[string]interface{}{"status": domain.StatusDelivered}),
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "not valid format id",
		},
		{
			name:        "invalid json",
			idParam:     orderID.String(),
			body:        `{"status":`,
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "unexpected EOF",
		},
		{
			name:        "missing status",
			idParam:     orderID.String(),
			body:        `{}`,
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Error:Field validation",
		},
		{
			name:    "service returns validation error",
			idParam: orderID.String(),
			body:    jsonBody(t, map[string]interface{}{"status": "processing"}),
			service: &fakeService{
				UpdateOrderStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
					return domain.ErrInvalidStatus
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: domain.ErrInvalidStatus.Error(),
		},
		{
			name:    "service returns not found",
			idParam: orderID.String(),
			body:    jsonBody(t, map[string]interface{}{"status": domain.StatusCancelled}),
			service: &fakeService{
				UpdateOrderStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
					return domain.ErrOrderNotFound
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: domain.ErrOrderNotFound.Error(),
		},
		{
			name:    "service returns unexpected error",
			idParam: orderID.String(),
			body:    jsonBody(t, map[string]interface{}{"status": domain.StatusCancelled}),
			service: &fakeService{
				UpdateOrderStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
					return errors.New("database update failed")
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "database update failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.service)

			w := performRequest(router, http.MethodPatch, "/orders/"+tt.idParam+"/status", tt.body)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.assertResponse != nil {
				tt.assertResponse(t, w.Body.Bytes())
				return
			}
			assertErrorMessageContains(t, w.Body.Bytes(), tt.wantMessage)
		})
	}
}

func TestOrderHandler_CancelOrder(t *testing.T) {
	orderID := uuid.New()

	tests := []struct {
		name           string
		idParam        string
		body           string
		service        *fakeService
		wantStatus     int
		wantMessage    string
		assertResponse func(t *testing.T, body []byte)
	}{
		{
			name:    "success",
			idParam: orderID.String(),
			body: jsonBody(t, map[string]interface{}{
				"reason": "customer changed their mind",
			}),
			service: &fakeService{
				CancelOrderFn: func(ctx context.Context, id uuid.UUID, reason string) error {
					assert.Equal(t, orderID, id)
					assert.Equal(t, "customer changed their mind", reason)
					return nil
				},
			},
			wantStatus: http.StatusOK,
			assertResponse: func(t *testing.T, body []byte) {
				var got map[string]string
				require.NoError(t, json.Unmarshal(body, &got))
				assert.Equal(t, domain.StatusCancelled, got["status"])
			},
		},
		{
			name:        "invalid uuid",
			idParam:     "not-a-uuid",
			body:        jsonBody(t, map[string]interface{}{"reason": "wrong order"}),
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "not valid format id",
		},
		{
			name:        "invalid json",
			idParam:     orderID.String(),
			body:        `{"reason":`,
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "unexpected EOF",
		},
		{
			name:        "missing reason",
			idParam:     orderID.String(),
			body:        `{}`,
			service:     &fakeService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Error:Field validation",
		},
		{
			name:    "service returns not found",
			idParam: orderID.String(),
			body:    jsonBody(t, map[string]interface{}{"reason": "wrong order"}),
			service: &fakeService{
				CancelOrderFn: func(ctx context.Context, id uuid.UUID, reason string) error {
					return domain.ErrOrderNotFound
				},
			},
			wantStatus:  http.StatusNotFound,
			wantMessage: "order not found",
		},
		{
			name:    "service returns cannot cancel",
			idParam: orderID.String(),
			body:    jsonBody(t, map[string]interface{}{"reason": "too late"}),
			service: &fakeService{
				CancelOrderFn: func(ctx context.Context, id uuid.UUID, reason string) error {
					return domain.ErrCannotCancelOrder
				},
			},
			wantStatus:  http.StatusConflict,
			wantMessage: domain.ErrCannotCancelOrder.Error(),
		},
		{
			name:    "service returns unexpected error",
			idParam: orderID.String(),
			body:    jsonBody(t, map[string]interface{}{"reason": "payment failed"}),
			service: &fakeService{
				CancelOrderFn: func(ctx context.Context, id uuid.UUID, reason string) error {
					return errors.New("database update failed")
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "database update failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.service)

			w := performRequest(router, http.MethodPatch, "/orders/"+tt.idParam+"/cancel", tt.body)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.assertResponse != nil {
				tt.assertResponse(t, w.Body.Bytes())
				return
			}
			assertErrorMessageContains(t, w.Body.Bytes(), tt.wantMessage)
		})
	}
}

func jsonBody(t *testing.T, payload interface{}) string {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	return string(body)
}

func performRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	if body == "" {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader([]byte(body))
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func assertErrorMessageContains(t *testing.T, body []byte, want string) {
	t.Helper()

	var got map[string]string
	require.NoError(t, json.Unmarshal(body, &got))
	assert.Contains(t, got["message"], want)
}

package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-service/internal/domain"

"github.com/gin-gonic/gin"
"github.com/google/uuid"
"github.com/stretchr/testify/assert"
)

// fakeService implements domain.OrderService for handler tests.
type fakeService struct {
    CreateOrderFn       func(ctx context.Context, customerID string, amount float64, status string) (*domain.Order, error)
    GetOrderFn          func(ctx context.Context, id uuid.UUID) (*domain.Order, error)
    UpdateOrderStatusFn func(ctx context.Context, id uuid.UUID, status string) error
}

func (f *fakeService) CreateOrder(ctx context.Context, customerID string, amount float64, status string) (*domain.Order, error) {
    if f.CreateOrderFn == nil {
        return nil, nil
    }
    return f.CreateOrderFn(ctx, customerID, amount, status)
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

func setupRouter(svc domain.OrderService) *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    h := NewOrderHandler(svc)
    r.POST("/orders", h.CreateOrder)
    return r
}

func TestOrderHandler_CreateOrder_Success(t *testing.T) {
    id := uuid.New()
    svc := &fakeService{
        CreateOrderFn: func(ctx context.Context, customerID string, amount float64, status string) (*domain.Order, error) {
            return &domain.Order{ID: id, CustomerID: uuid.MustParse(customerID), Amount: amount, Status: status}, nil
        },
    }

    router := setupRouter(svc)

    payload := map[string]interface{}{"customer_id": uuid.New().String(), "amount": 12.5, "status": "pending"}
    b, _ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    var got domain.Order
    err := json.Unmarshal(w.Body.Bytes(), &got)
    assert.NoError(t, err)
    assert.Equal(t, id.String(), got.ID.String())
    assert.Equal(t, "pending", got.Status)
}

func TestOrderHandler_CreateOrder_BadRequest(t *testing.T) {
    // Missing required fields -> binding error
    svc := &fakeService{}
    router := setupRouter(svc)

    payload := map[string]interface{}{"customer_id": "", "amount": 0}
    b, _ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    var resp map[string]interface{}
    _ = json.Unmarshal(w.Body.Bytes(), &resp)
    // helpers return message field
    assert.Contains(t, resp, "message")
}

func TestOrderHandler_CreateOrder_ServiceError(t *testing.T) {
    svc := &fakeService{
        CreateOrderFn: func(ctx context.Context, customerID string, amount float64, status string) (*domain.Order, error) {
            return nil, errors.New("service failure")
        },
    }
    router := setupRouter(svc)

    payload := map[string]interface{}{"customer_id": uuid.New().String(), "amount": 12.5, "status": "pending"}
    b, _ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    var resp map[string]interface{}
    _ = json.Unmarshal(w.Body.Bytes(), &resp)
    assert.Equal(t, "service failure", resp["message"])
}


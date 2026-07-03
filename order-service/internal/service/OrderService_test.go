package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"order-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// fakeRepo is a test double implementing domain.OrderRepository with injectable behavior.
type fakeRepo struct {
	CreateFn       func(ctx context.Context, order *domain.Order) error
	GetByIDFn      func(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	UpdateStatusFn func(ctx context.Context, id uuid.UUID, status string) error
}

func (f *fakeRepo) Create(ctx context.Context, order *domain.Order) error {
	if f.CreateFn == nil {
		return nil
	}
	return f.CreateFn(ctx, order)
}

func (f *fakeRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	if f.GetByIDFn == nil {
		return nil, nil
	}
	return f.GetByIDFn(ctx, id)
}

func (f *fakeRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	if f.UpdateStatusFn == nil {
		return nil
	}
	return f.UpdateStatusFn(ctx, id, status)
}

// TestOrderService_GetOrder tests the GetOrder method with various scenarios
func TestOrderService_GetOrder(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	tests := []struct {
		name      string
		repo      domain.OrderRepository
		wantErr   bool
		wantOrder *domain.Order
	}{
		{
			name: "success - order found",
			repo: &fakeRepo{
				GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
					return &domain.Order{ID: id, Status: "pending", Amount: 10.0, CreatedAt: time.Now()}, nil
				},
			},
			wantErr:   false,
			wantOrder: &domain.Order{ID: id},
		},
		{
			name: "not found - repo returns nil",
			repo: &fakeRepo{
				GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
					return nil, nil
				},
			},
			wantErr: true,
		},
		{
			name: "repo error",
			repo: &fakeRepo{
				GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.repo)
			got, err := svc.GetOrder(ctx, id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantOrder.ID, got.ID)
			}
		})
	}
}

// TestOrderService_CreateOrder tests the CreateOrder method with comprehensive validation cases
func TestOrderService_CreateOrder(t *testing.T) {
	ctx := context.Background()
	validCustomerID := uuid.New().String()

	tests := []struct {
		name       string
		repo       domain.OrderRepository
		customerID string
		amount     float64
		status     string
		wantErr    bool
		errMsg     string
		wantOrder  bool
	}{
		{
			name: "success - valid order creation",
			repo: &fakeRepo{
				CreateFn: func(ctx context.Context, order *domain.Order) error {
					assert.NotNil(t, order)
					assert.NotEqual(t, uuid.Nil, order.ID)
					assert.NotZero(t, order.CreatedAt)
					return nil
				},
			},
			customerID: validCustomerID,
			amount:     99.99,
			status:     "pending",
			wantErr:    false,
			wantOrder:  true,
		},
		{
			name:       "validation - invalid customer ID format",
			repo:       &fakeRepo{},
			customerID: "not-a-uuid",
			amount:     50.0,
			status:     "pending",
			wantErr:    true,
			errMsg:     "invalid customer id",
			wantOrder:  false,
		},
		{
			name:       "validation - negative amount",
			repo:       &fakeRepo{},
			customerID: validCustomerID,
			amount:     -10.0,
			status:     "pending",
			wantErr:    true,
			errMsg:     "amount must be greater than zero",
			wantOrder:  false,
		},
		{
			name:       "validation - zero amount",
			repo:       &fakeRepo{},
			customerID: validCustomerID,
			amount:     0,
			status:     "pending",
			wantErr:    true,
			errMsg:     "amount must be greater than zero",
			wantOrder:  false,
		},
		{
			name:       "validation - empty status",
			repo:       &fakeRepo{},
			customerID: validCustomerID,
			amount:     50.0,
			status:     "",
			wantErr:    true,
			errMsg:     "status is empty",
			wantOrder:  false,
		},
		{
			name: "error - repo create fails",
			repo: &fakeRepo{
				CreateFn: func(ctx context.Context, order *domain.Order) error {
					return errors.New("database error")
				},
			},
			customerID: validCustomerID,
			amount:     50.0,
			status:     "pending",
			wantErr:    true,
			wantOrder:  false,
		},
		{
			name: "success - shipped status",
			repo: &fakeRepo{
				CreateFn: func(ctx context.Context, order *domain.Order) error {
					return nil
				},
			},
			customerID: validCustomerID,
			amount:     25.50,
			status:     "shipped",
			wantErr:    false,
			wantOrder:  true,
		},
		{
			name: "success - large amount",
			repo: &fakeRepo{
				CreateFn: func(ctx context.Context, order *domain.Order) error {
					return nil
				},
			},
			customerID: validCustomerID,
			amount:     999999.99,
			status:     "pending",
			wantErr:    false,
			wantOrder:  true,
		},
		{
			name:       "validation - invalid UUID format",
			repo:       &fakeRepo{},
			customerID: "12345",
			amount:     100.0,
			status:     "pending",
			wantErr:    true,
			errMsg:     "invalid customer id",
			wantOrder:  false,
		},
		{
			name: "success - decimal amount",
			repo: &fakeRepo{
				CreateFn: func(ctx context.Context, order *domain.Order) error {
					assert.Equal(t, 0.01, order.Amount)
					return nil
				},
			},
			customerID: validCustomerID,
			amount:     0.01,
			status:     "pending",
			wantErr:    false,
			wantOrder:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.repo)
			// prepare domain.Order according to test case
			var cid uuid.UUID
			if parsed, perr := uuid.Parse(tt.customerID); perr == nil {
				cid = parsed
			} else {
				cid = uuid.Nil
			}
			orderArg := &domain.Order{CustomerID: cid, Amount: tt.amount, Status: tt.status}
			order, err := svc.CreateOrder(ctx, orderArg)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				if tt.wantOrder {
					assert.NotNil(t, order)
					assert.NotEqual(t, uuid.Nil, order.ID)
					assert.NotZero(t, order.CreatedAt)
				}
			}
		})
	}
}

// TestOrderService_UpdateOrderStatus tests status update with validation
func TestOrderService_UpdateOrderStatus(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	tests := []struct {
		name    string
		repo    domain.OrderRepository
		id      uuid.UUID
		status  string
		wantErr bool
		errMsg  string
	}{
		{
			name: "success - update to pending",
			repo: &fakeRepo{
				UpdateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
					return nil
				},
			},
			id:      id,
			status:  "pending",
			wantErr: false,
		},
		{
			name: "success - update to shipped",
			repo: &fakeRepo{
				UpdateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
					return nil
				},
			},
			id:      id,
			status:  "shipped",
			wantErr: false,
		},
		{
			name: "success - update to delivered",
			repo: &fakeRepo{
				UpdateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
					return nil
				},
			},
			id:      id,
			status:  "delivered",
			wantErr: false,
		},
		{
			name: "success - update to cancelled",
			repo: &fakeRepo{
				UpdateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
					return nil
				},
			},
			id:      id,
			status:  "cancelled",
			wantErr: false,
		},
		{
			name:    "validation - empty status",
			repo:    &fakeRepo{},
			id:      id,
			status:  "",
			wantErr: true,
			errMsg:  "status is empty",
		},
		{
			name:    "validation - invalid status unknown",
			repo:    &fakeRepo{},
			id:      id,
			status:  "unknown",
			wantErr: true,
			errMsg:  "status is invalid",
		},
		{
			name:    "validation - invalid status paid",
			repo:    &fakeRepo{},
			id:      id,
			status:  "paid",
			wantErr: true,
			errMsg:  "status is invalid",
		},
		{
			name:    "validation - invalid status processing",
			repo:    &fakeRepo{},
			id:      id,
			status:  "processing",
			wantErr: true,
			errMsg:  "status is invalid",
		},
		{
			name: "error - repo update fails",
			repo: &fakeRepo{
				UpdateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
					return errors.New("database update failed")
				},
			},
			id:      id,
			status:  "shipped",
			wantErr: true,
		},
		{
			name:    "validation - case sensitive status",
			repo:    &fakeRepo{},
			id:      id,
			status:  "Pending",
			wantErr: true,
			errMsg:  "status is invalid",
		},
		{
			name:    "validation - invalid status with spaces",
			repo:    &fakeRepo{},
			id:      id,
			status:  " pending ",
			wantErr: true,
			errMsg:  "status is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.repo)
			err := svc.UpdateOrderStatus(ctx, tt.id, tt.status)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}


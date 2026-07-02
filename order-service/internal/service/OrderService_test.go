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
            name: "success",
            repo: &fakeRepo{
                GetByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
                    return &domain.Order{ID: id, Status: "pending", Amount: 10.0, CreatedAt: time.Now()}, nil
                },
            },
            wantErr:   false,
            wantOrder: &domain.Order{ID: id},
        },
        {
            name: "not found (nil, nil)",
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
            svc := newOrderService(tt.repo)
            got, err := svc.GetOrder(ctx, id)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantOrder.ID, got.ID)
            }
        })
    }
}

func TestOrderService_PlaceOrder(t *testing.T) {
    ctx := context.Background()
    now := time.Now()

    tests := []struct {
        name    string
        repo    domain.OrderRepository
        order   *domain.Order
        wantErr bool
    }{
        {
            name: "success create",
            repo: &fakeRepo{
                CreateFn: func(ctx context.Context, order *domain.Order) error { return nil },
            },
            order: &domain.Order{ID: uuid.New(), CustomerID: uuid.New(), Status: "pending", Amount: 5, CreatedAt: now},
            wantErr: false,
        },
        {
            name: "nil order",
            repo: &fakeRepo{},
            order: nil,
            wantErr: true,
        },
        {
            name: "repo create error",
            repo: &fakeRepo{
                CreateFn: func(ctx context.Context, order *domain.Order) error { return errors.New("create failed") },
            },
            order: &domain.Order{ID: uuid.New()},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := newOrderService(tt.repo)
            err := svc.PlaceOrder(ctx, tt.order)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestOrderService_UpdateOrderStatus(t *testing.T) {
    ctx := context.Background()
    id := uuid.New()

    tests := []struct {
        name      string
        repo      domain.OrderRepository
        id        uuid.UUID
        status    string
        wantErr   bool
    }{
        {
            name: "success update",
            repo: &fakeRepo{UpdateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error { return nil }},
            id: id,
            status: "shipped",
            wantErr: false,
        },
        {
            name: "empty status",
            repo: &fakeRepo{},
            id: id,
            status: "",
            wantErr: true,
        },
        {
            name: "invalid status",
            repo: &fakeRepo{},
            id: id,
            status: "unknown",
            wantErr: true,
        },
        {
            name: "repo error",
            repo: &fakeRepo{UpdateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error { return errors.New("update failed") }},
            id: id,
            status: "paid",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := newOrderService(tt.repo)
            err := svc.UpdateOrderStatus(ctx, tt.id, tt.status)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}


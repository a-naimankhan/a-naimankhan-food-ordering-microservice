package repository

import (
	"context"
	"database/sql"
	"order-service/internal/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	var repo domain.OrderRepository = NewOrderRepo(sqlxDB)

	tests := []struct {
		name    string
		order   *domain.Order
		mockErr error
		wantErr bool
	}{
		{
			name: "Success",
			order: &domain.Order{
				ID:         uuid.New(),
				CustomerID: uuid.New(),
				Status:     "pending",
				Amount:     99.9,
				CreatedAt:  time.Now(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "DB Error",
			order: &domain.Order{
				ID: uuid.New(),
			},
			mockErr: assert.AnError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mock.ExpectExec("INSERT INTO ORDERS").
					WithArgs(tt.order.ID, tt.order.CustomerID, tt.order.Status, tt.order.Amount, tt.order.CreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("INSERT INTO ORDERS").
					WillReturnError(tt.mockErr)
			}

			err := repo.Create(context.Background(), tt.order)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderRepo_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	var repo domain.OrderRepository = NewOrderRepo(sqlxDB)

	orderID := uuid.New()
	customerID := uuid.New()

	tests := []struct {
		name          string
		id            uuid.UUID
		mockSetup     func()
		wantErr       bool
		expectedOrder *domain.Order
	}{
		{
			name: "Success",
			id:   orderID,
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM ORDERS").
					WithArgs(orderID). // У тебя в SQL сейчас $1, а не $2, поправь запрос!
					WillReturnRows(sqlmock.NewRows([]string{"id", "customer_id", "status", "amount", "created_at"}).
						AddRow(orderID, customerID, "pending", 99.9, time.Now()))
			},
			wantErr:       false,
			expectedOrder: &domain.Order{ID: orderID},
		},
		{
			name: "Not Found",
			id:   uuid.New(),
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM ORDERS").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			order, err := repo.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOrder.ID, order.ID)
			}
		})
	}
}

func TestOrderRepo_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	var repo domain.OrderRepository = NewOrderRepo(sqlxDB)

	orderID := uuid.New()
	newStatus := "paid"

	tests := []struct {
		name          string
		id            uuid.UUID
		mockSetup     func()
		wantErr       bool
		expectedOrder *domain.Order
	}{
		{
			name: "SUCCESS",
			id:   orderID,
			mockSetup: func() {
				mock.ExpectExec("UPDATE ORDERS").
					WithArgs(orderID, newStatus).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr:       false,
			expectedOrder: &domain.Order{ID: orderID, Status: newStatus},
		},
		{
			name: "DB Error",
			id:   uuid.New(),
			mockSetup: func() {
				mock.ExpectExec("UPDATE ORDERS").
					WithArgs(orderID, newStatus).
					WillReturnError(assert.AnError)
			},
			wantErr:       true,
			expectedOrder: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.UpdateStatus(context.Background(), tt.id, newStatus)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

package domain

import "errors"

var (
    ErrOrderNotFound      = errors.New("order not found")
    ErrInvalidAmount      = errors.New("amount must be greater than zero")
    ErrInvalidCustomerID  = errors.New("invalid customer id")
    ErrStatusEmpty        = errors.New("status is empty")
    ErrInvalidStatus      = errors.New("status is invalid")
)



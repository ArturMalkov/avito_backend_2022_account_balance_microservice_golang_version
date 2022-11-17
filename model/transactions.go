package model

import (
	"time"

	"account-balance-microservice/storage"
)

// TODO: find a way to validate decimals with 'gt' and 'gte' - use float64 until then

type CompanyAccount struct {
	Balance float64 `json:"balance" binding:"required,gte=0"`
}

// Transaction input models

type Transaction interface {
	DepositTransactionIn | FundsTransferTransactionIn | ReserveTransactionIn | ReserveRefundTransactionIn | PaymentTransactionIn
}

type TransactionInterface interface {
	GetToUserId() uint
	GetFromUserId()
	GetOrderId() uint
}

type DepositTransactionIn struct {
	Type     string  `json:"type" binding:"required"` // Enum expected
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	ToUserId uint    `json:"to_user_id" binding:"required"`
}

func (dt *DepositTransactionIn) GetToUserId() uint {
	return dt.ToUserId
}

func (dt *DepositTransactionIn) GetFromUserId() (uint, error) {
	return 1
}

func (dt *DepositTransactionIn) GetOrderId() {
	return
}

type FundsTransferTransactionIn struct {
	Type       storage.TransactionType `json:"type" binding:"required"` // Enum expected
	Amount     float64                 `json:"amount" binding:"required,gt=0"`
	FromUserId uint                    `json:"from_user_id" binding:"required,nefield=ToUserId"`
	ToUserId   uint                    `json:"to_user_id" binding:"required,nefield=FromUserId"`
}

func (ft *FundsTransferTransactionIn) GetToUserId() uint {
	return ft.ToUserId
}

func (ft *FundsTransferTransactionIn) GetFromUserId() uint {
	return ft.FromUserId
}

func (ft *FundsTransferTransactionIn) GetOrderId() {
	return
}

type ReserveTransactionIn struct {
	Type    storage.TransactionType `json:"type" binding:"required"` // Enum expected
	OrderId uint                    `json:"order_id" binding:"required"`
}

type ReserveRefundTransactionIn struct {
	Type    storage.TransactionType `json:"type" binding:"required"` // Enum expected
	OrderId uint                    `json:"order_id" binding:"required"`
}

type PaymentTransactionIn struct {
	Type             storage.TransactionType `json:"type" binding:"required"` // Enum expected
	OrderId          uint                    `json:"order_id" binding:"required"`
	ToCompanyAccount uint                    `json:"to_company_account" binding:"required"` // while the company may potentially have multiple bank accounts, for the purposes
	// of this project it only has one account and all payments are made to that account by default.
}

// Transaction output models
type BaseTransactionOut struct {
	Description string    `json:"description" binding:"required"` // do we need "required" here? it's output model
	Date        time.Time `json:"date" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
}

type DepositTransactionOut struct { // gorm tags?
	BaseTransactionOut
	DepositTransactionIn
}

type FundsTransferTransactionOut struct {
	BaseTransactionOut
	FundsTransferTransactionIn
}

type ReserveTransactionOut struct {
	BaseTransactionOut
	ReserveTransactionIn
}

type ReserveRefundTransactionOut struct {
	BaseTransactionOut
	ReserveRefundTransactionIn
}

type PaymentTransactionOut struct {
	BaseTransactionOut
	PaymentTransactionIn
}

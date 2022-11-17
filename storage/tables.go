package storage

import (
	"time"
	// "github.com/shopspring/decimal"
)

type User struct {
	Id           uint   `gorm:"primaryKey"`
	FirstName    string `gorm:"size:50"`
	LastName     string `gorm:"size:50"`
	Username     string `gorm:"size:50; unique; not null"`
	Email        string `gorm:"size:50; unique; not null"`
	PhoneNumber  string `gorm:"size:50; unique; not null"`
	UserAccounts [2]UserAccount
	Orders       []Order
}

// find a way to use Enums with gorm models - now values columns are simply of string type

type AccountType string

const (
	REGULAR AccountType = "REGULAR"
	RESERVE AccountType = "RESERVE"
)

// find a way to deal with decimals - use float64 for now

type Account interface {
	IncreaseBalance(amount float64)
}

type UserAccount struct {
	Id      uint        `gorm:"primaryKey"`
	Type    AccountType `gorm:"not null"`
	Balance float64     `gorm:"default:0"`
	UserId  uint        `gorm:"index"`
}

func (u *UserAccount) IncreaseBalance(amount float64) {
	u.Balance += amount
}

type CompanyAccount struct {
	Id                uint    `gorm:"primaryKey"`
	Balance           float64 `gorm:"default:0"`
	BankAccountNumber string  `gorm:"size:50; not null"`
	Bank              string  `gorm:"size:50; not null"`
}

func (c *CompanyAccount) IncreaseBalance(amount float64) {
	c.Balance += amount
}

// find a way to use Enums with gorm models - now values columns are simply of string type

type OrderStatus string

const (
	NOT_SUBMITTED OrderStatus = "NOT_SUBMITTED"
	IN_PROGRESS   OrderStatus = "IN_PROGRESS"
	COMPLETED     OrderStatus = "COMPLETED"
	CANCELLED     OrderStatus = "CANCELLED"
)

type Order struct {
	Id           uint        `gorm:"primaryKey"`
	Status       OrderStatus `gorm:"default:'NOT_SUBMITTED'"`
	UserId       uint        `gorm:"index"`
	Items        []OrderItem
	Transactions []Transaction
}

type OrderItem struct {
	Id        uint `gorm:"primaryKey"`
	Quantity  uint `gorm:"default:1"`
	ServiceId uint `gorm:"index"`
	OrderId   uint `gorm:"index"`
}

type Service struct {
	Id              uint    `gorm:"primaryKey"`
	Name            string  `gorm:"size:50; not null"`
	Price           float64 `gorm:"not null"`
	Description     string  `gorm:"size:255; not null"`
	OrderedServices []OrderItem
}

// find a way to use Enums with gorm models - now values columns are simply of string type

type TransactionType string

const (
	DEPOSIT            TransactionType = "DEPOSIT"
	FUNDS_TRANSFER     TransactionType = "FUNDS_TRANSFER"
	RESERVE_FUNDS      TransactionType = "RESERVE_FUNDS"
	RESERVE_REFUND     TransactionType = "RESERVE_REFUND"
	PAYMENT_TO_COMPANY TransactionType = "PAYMENT_TO_COMPANY"
)

type Transaction struct {
	Id               uint `gorm:"primaryKey"`
	Amount           float64
	Type             TransactionType `gorm:"not null"`
	Description      string          `gorm:"size:255; not null"`
	Date             time.Time       `gorm:"default:1"` // default=datetime.datetime.utcnow
	OrderId          uint            `gorm:"index; default: NULL"`
	FromUserId       uint            `gorm:"default: NULL"`
	ToUserId         uint            `gorm:"default: NULL"`
	ToCompanyAccount uint            `gorm:"default: NULL"`
}

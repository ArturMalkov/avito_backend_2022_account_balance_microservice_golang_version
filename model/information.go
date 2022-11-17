package model

import (
	"account-balance-microservice/storage"
)

// https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go#L18

type UserAccountOut struct {
	UserId  uint                `json:"user_id" validate:"required" binding:"required"` // do we need required here??? it's output model
	Type    storage.AccountType `json:"type" validate:"required" binding:"required"`    // Enum member here!
	Balance float64             `json:"balance" validate:"required, gte=0" binding:"required"`
}

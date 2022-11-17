package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"account-balance-microservice/model"
	"account-balance-microservice/service"
)

// Deposits a specified amount of money to a specific user (i.e. user account balance increases) via external services
// (which are out of scope for this project).
// If user doesn't have an account yet, account will be created (as per the project's requirements)
// and money will be deposited to the account.
func depositFundsToAccount(transactionsService *service.TransactionsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var depositTransactionData model.DepositTransactionIn
		err := ctx.ShouldBindJSON(&depositTransactionData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()}) // or 422 - unprocessable entity?
			return
		}

		results := transactionsService.DepositFundsToAccount(depositTransactionData) // model.DepositTransactionOut
		ctx.JSON(http.StatusOK, results)
	}
}

// Transfers a specified amount of money from one user to another (i.e. sender account balance decreases while
// recipient account balance increases accordingly).
// If a recipient user doesn't have an account yet, account will be created (as per the project's requirements)
// and money will be transferred to the account.
func transferFundsBetweenUserAccounts(transactionsService *service.TransactionsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var fundsTransferTransactionData model.FundsTransferTransactionIn
		err := ctx.Bind(&fundsTransferTransactionData)
		if err != nil {
			panic(err) // validation error
		}

		results := transactionsService.TransferFundsBetweenUserAccounts(fundsTransferTransactionData) // model.FundsTransferTransactionOut
		ctx.JSON(http.StatusOK, results)
	}
}

// Reserves money from account of a user which made a specific order (money is transferred from user's regular account
// to reserve one).
// Changes the order's status to "in progress".
// Reserved money can later be refunded to regular account (if the order is cancelled)
// or paid to the company (if the order is fulfilled).
// The amount of money to be reserved is determined by the total price of the services in the order.
func reserveFunds(transactionsService *service.TransactionsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reserveTransactionData model.ReserveTransactionIn
		err := ctx.ShouldBindJSON(&reserveTransactionData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
			return
		}

		results := transactionsService.ReserveFunds(reserveTransactionData) // model.ReserveTransactionOut
		ctx.JSON(http.StatusOK, results)
	}
}

// Refunds previously reserved money (money is transferred back from user's reserve account to regular one) as per
// specific order.
// Changes the order's status to "cancelled".
func cancelReserve(transactionsService *service.TransactionsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reserveRefundTransactionData model.ReserveRefundTransactionIn
		err := ctx.ShouldBindJSON(&reserveRefundTransactionData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
			return
		}

		results := transactionsService.CancelReserve(reserveRefundTransactionData) // model.ReserveRefundTransactionOut
		ctx.JSON(http.StatusOK, results)
	}
}

// Transfers previously reserved money to company account (money is transferred from user's reserve account to
// company account).
// Changes the order's status to "completed".
func makePaymentToCompany(transactionsService *service.TransactionsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var paymentTransactionData model.PaymentTransactionIn
		err := ctx.ShouldBindJSON(&paymentTransactionData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
			return
		}

		results := transactionsService.MakePaymentToCompany(paymentTransactionData) // model.PaymentTransactionOut
		ctx.JSON(http.StatusOK, results)
	}
}

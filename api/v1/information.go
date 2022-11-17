package api

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"account-balance-microservice/service"

	"strconv"
)

func getAccountBalanceInfo(informationService *service.InformationService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, _ := strconv.Atoi(ctx.Param("user_id"))

		results := informationService.GetAccountBalanceInfo(uint(userId)) // []model.UserAccountOut
		ctx.JSON(http.StatusOK, results)
	}
}

func getAccountTransactionsInfo(informationService *service.InformationService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, _ := strconv.Atoi(ctx.Param("user_id"))
		page, _ := strconv.Atoi(ctx.Query("page"))
		sortByAmount, _ := strconv.ParseBool(ctx.Query("sort_by_amount"))
		sortByDate, _ := strconv.ParseBool(ctx.Query("sort_by_date"))

		results := informationService.GetAccountTransactionsInfo(uint(userId), page, sortByAmount, sortByDate) // []model.DepositTransactionOut, model.FundsTransferTransactionOut, model.ReserveTransactionOut, model.ReserveRefundTransactionOut, model.PaymentTransactionOut
		ctx.JSON(http.StatusOK, results)
	}
}

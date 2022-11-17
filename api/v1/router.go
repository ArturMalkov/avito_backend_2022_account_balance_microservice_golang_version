package api

import (
	"account-balance-microservice/service"
	"account-balance-microservice/storage"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	databaseConnection := *storage.GetDatabaseConnection()

	informationService := service.InformationService{&databaseConnection}
	reportsService := service.ReportsService{&databaseConnection}
	transactionsService := service.TransactionsService{&databaseConnection, &informationService}

	router := gin.Default()
	apiV1Router := router.Group("/v1")

	informationRouter := apiV1Router.Group("/information")
	{
		informationRouter.GET("/account-balance/:user_id", getAccountBalanceInfo(&informationService))
		informationRouter.GET("/account-transactions/:user_id", getAccountTransactionsInfo(&informationService))
	}

	reportsRouter := apiV1Router.Group("/reports")
	{
		reportsRouter.GET("/consolidated/monthly", getMonthlyAccountingReport(&reportsService))
	}

	transactionsRouter := apiV1Router.Group("/transactions")
	{
		transactionsRouter.PATCH("/deposit", depositFundsToAccount(&transactionsService))
		transactionsRouter.PATCH("/transfer", transferFundsBetweenUserAccounts(&transactionsService))
		transactionsRouter.PATCH("/reserve", reserveFunds(&transactionsService))
		transactionsRouter.PATCH("/reserve-refund", cancelReserve(&transactionsService))
		transactionsRouter.PATCH("/make-payment", makePaymentToCompany(&transactionsService))
	}

	return router
}

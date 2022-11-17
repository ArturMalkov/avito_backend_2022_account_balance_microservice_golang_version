package api

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"strconv"

	"account-balance-microservice/service"
)

// Returns csv report with total revenues for each service rendered in the requested period.
// Format: service name, total revenues in the reporting period.
func getMonthlyAccountingReport(reportsService *service.ReportsService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		year, _ := strconv.Atoi(ctx.Query("year"))
		month, _ := strconv.Atoi(ctx.Query("page"))

		report := reportsService.PrepareMonthlyAccountingReportInCSV(year, month)
		ctx.JSON(http.StatusOK, report)
	}
}

package service

import (
	"gorm.io/gorm"
)

type ReportsService struct {
	Db *gorm.DB
}

func (service ReportsService) PrepareMonthlyAccountingReportInCSV(year int, month int) int {
	return 1
}

// Calculates revenue for each service rendered in the reporting period.
// Returns a dictionary which maps service name with total revenue from it in the period.
func (service ReportsService) calculateRevenueFromServices(year int, month int) {

}

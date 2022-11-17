package service

import (
	"account-balance-microservice/storage"

	"github.com/oleiade/reflections"

	"gorm.io/gorm"
)

// Responsible for working with user accounts' data.
type InformationService struct {
	Db *gorm.DB
}

// Returns info on user's regular and reserve accounts.
func (service InformationService) GetAccountBalanceInfo(userId uint) [2]storage.UserAccount { // TODO: convert table query results to models
	regularAccountBalanceInfo := service.getUserAccountByUserId(userId, "REGULAR")
	reserveAccountBalanceInfo := service.getUserAccountByUserId(userId, "RESERVE")

	accountsInfo := [2]storage.UserAccount{regularAccountBalanceInfo, reserveAccountBalanceInfo}
	return accountsInfo
}

// Returns a list of transactions with description on where and why the funds
// were credited/debited from the account balance.
// Sorting (by date and amount) and pagination of results are provided as an option.
func (service InformationService) GetAccountTransactionsInfo(userId uint, page int, sortByAmount bool, sortByDate bool) int {
	return 1
}

// Returns account of a specified type (regular/reserve) of a particular user.
func (service InformationService) getUserAccountByUserId(userId uint, accountType string) storage.UserAccount {
	service.RaiseErrorIfUserDoesNotExist(userId)

	var userAccount storage.UserAccount
	if err := service.Db.Where(storage.UserAccount{UserId: userId, Type: "REGULAR"}).First(&userAccount).Error; err != nil {
		panic("user account doesn't exist'")
	}

	return userAccount

}

// Returns both regular and reserve accounts of a particular user.
func (service InformationService) getUserAccountsByOrderId(orderId uint) (storage.UserAccount, storage.UserAccount) {
	var order storage.Order
	service.Db.First(&order, orderId)

	var user storage.User
	service.Db.First(&user, order.UserId)

	var userAccounts [2]storage.UserAccount
	userAccounts = user.UserAccounts

	if &userAccounts == nil {
		panic("some message here!")
	}

	regularAccount := userAccounts[0]
	reserveAccount := userAccounts[1]

	return regularAccount, reserveAccount
}

func (service InformationService) getCompanyAccountByCompanyAccountId(companyAccountId uint) storage.CompanyAccount {
	var companyAccount storage.CompanyAccount
	service.Db.First(&companyAccount, companyAccountId)

	if &companyAccount == nil {
		panic("Couldn't find company account'")
	}

	return companyAccount
}

func (service InformationService) getUserByUserId(userId uint) storage.User {
	service.RaiseErrorIfUserDoesNotExist(userId)

	var user storage.User
	service.Db.First(&user, userId)

	return user
}

func (service InformationService) RaiseErrorIfUserDoesNotExist(userId uint) {
	var user storage.User
	result := service.Db.First(&user, userId)

	if result == nil {
		panic("user does not exist")
	}
}

func getTablePaginationResults(selectedRows string, pageNumber int) {

}

// Only works with tables with 'date' field present.
func sortRowsByDateColumn(table any, selectedRows string) {
	has, _ := reflections.HasField(table, "Amount")
	if has == false {
		panic("Table does not have 'Amount' field")
	}

}

// Only works with tables with 'amount' field present.
func sortRowsByAmountColumn(table any, selectedRows string) {
	has, _ := reflections.HasField(table, "Amount")
	if has == false {
		panic("Table does not have 'Amount' field")
	}

}

package service

import (
	"gorm.io/gorm"

	// "net/http"

	"account-balance-microservice/model"

	"account-balance-microservice/storage"
)

type TransactionDesciption string

const (
	DEPOSIT            TransactionDesciption = "Money in the amount of %sUSD was deposited to user %s from external services on %s."
	FUNDS_TRANSFER     TransactionDesciption = "Money in the amount of %sUSD was transferred from user %s to user %s on %s."
	RESERVE_FUNDS      TransactionDesciption = "Money in the amount of %sUSD was reserved on user %s reserve account as per the order %s on %s."
	RESERVE_REFUND     TransactionDesciption = "Money in the amount of %sUSD was refunded to user %s account from his/her reserve account as per the order %s on %s."
	PAYMENT_TO_COMPANY TransactionDesciption = "Money in the amount of %sUSD was paid by user %s to the company account as per the order %s on %s."
)

// Responsible for working with transactions' and accounts' data.
type TransactionsService struct {
	Db          *gorm.DB
	InfoService *InformationService
}

func (service *TransactionsService) DepositFundsToAccount(transactionData model.DepositTransactionIn) storage.Transaction {
	recipientAccount := service.getOrCreateRecipientAccountByUserId(transactionData.ToUserId)

	transferFunds(
		nil,
		&recipientAccount,
		transactionData.Amount,
	)

	transaction := prepareTransactionDBRecord(
		transactionData,
		storage.DEPOSIT,
		transactionData.Amount,
		nil,
		&transactionData.ToUserId,
		nil,
	)

	// TODO: find a way to create all the records below in one transaction
	service.Db.Create(&recipientAccount)
	service.Db.Create(&transaction)

	return transaction
}

// Transfers funds between regular(!) accounts of two users.
func (service *TransactionsService) TransferFundsBetweenUserAccounts(transactionData model.FundsTransferTransactionIn) storage.Transaction {
	senderAccount := service.InfoService.getUserAccountByUserId(
		transactionData.FromUserId,
		"REGULAR",
	)
	recipientAccount := service.getOrCreateRecipientAccountByUserId(transactionData.ToUserId)

	transferFunds(
		&senderAccount,
		&recipientAccount,
		transactionData.Amount,
	)

	transaction := prepareTransactionDBRecord(
		transactionData,
		storage.FUNDS_TRANSFER,
		transactionData.Amount,
		&transactionData.FromUserId,
		&transactionData.ToUserId,
		nil,
	)

	// TODO: find a way to create all the records below in one transaction
	service.Db.Create(&senderAccount)
	service.Db.Create(&recipientAccount)
	service.Db.Create(&transaction)

	return transaction
}

func (service *TransactionsService) ReserveFunds(transactionData model.ReserveTransactionIn) storage.Transaction {
	service.raiseErrorIfOrderDoesNotExistOrHasIrrelevantStatus(
		transactionData.OrderId,
		storage.NOT_SUBMITTED,
	)

	regularAccount, reserveAccount := service.InfoService.getUserAccountsByOrderId(transactionData.OrderId)

	transactionAmount := service.calculateTransactionAmount(transactionData.OrderId)

	transferFunds(
		&regularAccount,
		&reserveAccount,
		transactionAmount,
	)

	transaction := prepareTransactionDBRecord(
		transactionData,
		storage.RESERVE_FUNDS,
		transactionAmount,
		&regularAccount.UserId,
		nil,
		&transactionData.OrderId,
	)

	// TODO: find a way to create all the records below in one transaction
	service.Db.Create(&regularAccount)
	service.Db.Create(&reserveAccount)
	service.Db.Create(&transaction)

	// TODO: find a way to update order status within the same transaction which creates the records above
	updateOrderStatus(
		transactionData.OrderId,
		storage.IN_PROGRESS,
		service.Db,
	)

	return transaction
}

func (service *TransactionsService) CancelReserve(transactionData model.ReserveRefundTransactionIn) storage.Transaction {
	service.raiseErrorIfOrderDoesNotExistOrHasIrrelevantStatus(
		transactionData.OrderId,
		storage.IN_PROGRESS,
	)

	reserveTransactionToBeCancelled := service.getTransactionByOrderId(
		transactionData.OrderId,
		"RESERVE",
	)

	transactionAmount := reserveTransactionToBeCancelled.Amount
	regularAccount, reserveAccount := service.InfoService.getUserAccountsByOrderId(
		transactionData.OrderId,
	)

	transferFunds(
		&reserveAccount,
		&regularAccount,
		transactionAmount,
	)

	transaction := prepareTransactionDBRecord(
		transactionData,
		storage.RESERVE_REFUND,
		transactionAmount,
		nil,
		&regularAccount.UserId,
		&transactionData.OrderId,
	)

	// TODO: find a way to create all the records below in one transaction
	service.Db.Create(&regularAccount)
	service.Db.Create(&reserveAccount)
	service.Db.Create(&transaction)

	// TODO: find a way to update order status within the same transaction which creates the records above
	updateOrderStatus(
		transactionData.OrderId,
		storage.CANCELLED,
		service.Db,
	)
	// self.db_session.commit()

	return transaction
}

// Transfers reserved (as per specified order) money from user's reserve account to company account.
// While the company can potentially have multiple bank accounts, for the purposes of this project it only has one
// account and all payments are made to that account by default.
func (service *TransactionsService) MakePaymentToCompany(transactionData model.PaymentTransactionIn) storage.Transaction {
	service.raiseErrorIfOrderDoesNotExistOrHasIrrelevantStatus(
		transactionData.OrderId,
		storage.IN_PROGRESS,
	)

	reserveTransactionToBePaid := service.getTransactionByOrderId(
		transactionData.OrderId, "RESERVE_FUNDS",
	)
	transactionAmount := reserveTransactionToBePaid.Amount
	_, reserveAccount := service.InfoService.getUserAccountsByOrderId(transactionData.OrderId)
	companyAccount := service.InfoService.getCompanyAccountByCompanyAccountId(transactionData.ToCompanyAccount)

	transferFunds(
		&reserveAccount,
		&companyAccount,
		transactionAmount,
	)

	transaction := prepareTransactionDBRecord(
		transactionData,
		storage.PAYMENT_TO_COMPANY,
		transactionAmount,
		&reserveAccount.UserId,
		nil,
		&transactionData.OrderId,
	)

	// TODO: find a way to create all the records below in one transaction
	service.Db.Create(&reserveAccount)
	service.Db.Create(&companyAccount)
	service.Db.Create(&transaction)

	// TODO: find a way to update order status within the same transaction which creates the records above
	updateOrderStatus(
		transactionData.OrderId,
		storage.COMPLETED,
		service.Db,
	)
	// self.db_session.commit()

	return transaction
}

func (service *TransactionsService) calculateTransactionAmount(orderId uint) float64 {
	// order = self.db_session.query(tables.Order).filter(tables.Order.id == order_id).first()
	// transaction_amount = Decimal(0)
	// for ordered_service in order.items:
	//     transaction_amount += ordered_service.quantity * ordered_service.service.price
	// return transaction_amount
	return 100.5
}

func (service *TransactionsService) getTransactionByOrderId(orderId uint, type_ storage.TransactionType) storage.Transaction {
	var transaction storage.Transaction
	if err := service.Db.Where(storage.Transaction{OrderId: orderId, Type: type_}).First(&transaction).Error; err != nil {
		// http.Error(message.TRANSACTION_DOES_NOT_EXIST, 422)
		panic("Transaction does not exist!")
	}

	return transaction

}

// In case of deposit/money transfer transactions, if a recipient user doesn't have an account yet,
// both regular and reserve accounts will be automatically created with zero balance.
func (service *TransactionsService) getOrCreateRecipientAccountByUserId(userId uint) storage.UserAccount {
	service.InfoService.RaiseErrorIfUserDoesNotExist(userId)

	var regularAccount storage.UserAccount
	if err := service.Db.Where(storage.UserAccount{UserId: userId, Type: "REGULAR"}).First(&regularAccount).Error; err != nil {
		regularAccount := storage.UserAccount{UserId: userId, Type: storage.REGULAR}
		reserveAccount := storage.UserAccount{UserId: userId, Type: storage.RESERVE}
		service.Db.Create(&regularAccount)
		service.Db.Create(&reserveAccount)
		// self.db_session.commit()
	}

	return regularAccount
}

// https://stackoverflow.com/questions/71376627/in-go-generics-how-to-use-a-common-method-for-types-in-a-union-constraint
func prepareTransactionDBRecord[T model.Transaction](transactionData T, type_ storage.TransactionType, amount float64, fromUserId *uint, toUserId *uint, orderId *uint) storage.Transaction {
	transaction := storage.Transaction{
		FromUserId:  transactionData.GetFromUserId(),
		ToUserId:    transactionData.GetToUserId(),
		OrderId:     transactionData.GetOrderId(),
		Type:        type_,
		Amount:      amount,
		Description: service.prepareTransactionDescription(type_, amount, *fromUserId, *toUserId, *orderId),
	}

	return transaction
}

// Utility method used to transfer funds in all transactions except for deposits.
// Do not confuse it with transfer_funds_between_user_accounts().
// func transferFunds[T *storage.UserAccount | *storage.CompanyAccount](senderAccount *storage.UserAccount, recipientAccount T, transferAmount float64) {
func transferFunds(senderAccount *storage.UserAccount, recipientAccount storage.Account, transferAmount float64) {

	if senderAccount == nil { // there's no sender account in deposit transaction
		senderAccount.Balance -= transferAmount
		if senderAccount.Balance < 0 {
			// http.Error(message.ACCOUNT_BALANCE_CANNOT_BE_NEGATIVE, 422)
			panic("Account cannot be negative!")
		}
	}
	recipientAccount.IncreaseBalance(transferAmount)
}

// TODO: rewrite this function from the Python version of the project
func (service *TransactionsService) prepareTransactionDescription(transactionType storage.TransactionType, amount float64, fromUserId uint, toUserId uint, orderId uint) string {
	return 1
}

// Helps to avoid situations when order doesn't exist or has irrelevant status with regard to the transaction
// to be performed.
// E.g.:
// for 'reserve' transaction to happen, order must be in 'not submitted' state;
// for 'reserve refund' or 'payment to company' transaction to happen, order must be in 'in progress' state.
func (service TransactionsService) raiseErrorIfOrderDoesNotExistOrHasIrrelevantStatus(orderId uint, correctStatus storage.OrderStatus) {
	var order storage.Order

	if err := service.Db.First(&order, orderId).Error; err != nil {
		// http.Error(message.ORDER_DOES_NOT_EXIST, 422)
		panic("Order does not exist!")
	}

	if order.Status != correctStatus {
		// http.Error(message.INCORRECT_ORDER_STATUS, 422)
		panic("Incorrect order status!")
	}
}

// This method should belong to another microservice - orders microservice.
// Placed here for demonstration/convenience purposes.
func updateOrderStatus(orderId uint, newStatus storage.OrderStatus, db *gorm.DB) {
	var order storage.Order

	db.First(&order, orderId)
	order.Status = newStatus
	db.Save(&order)
	// no commit here - commits happen in the methods of the TransactionsService
	// (account balance operations, transactions' saves and order status updates all need to happen within the same
	// transaction)
}

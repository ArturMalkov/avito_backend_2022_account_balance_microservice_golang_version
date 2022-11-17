package message

type ExceptionDescription string

const (
	ACCOUNT_BALANCE_CANNOT_BE_NEGATIVE             ExceptionDescription = "User account balance cannot be negative. Transaction failed."
	SENDER_AND_RECIPIENT_CANNOT_BE_THE_SAME_PERSON ExceptionDescription = "Sender and recipient cannot be the same person."
	ACCOUNT_DOES_NOT_EXIST                         ExceptionDescription = "User doesn't have an account yet. Please deposit or transfer money to that user."
	USER_DOES_NOT_EXIST                            ExceptionDescription = "Specified user doesn't exist."
	RESOURCE_CANNOT_BE_SORTED_BY_DATE              ExceptionDescription = "Requested resource cannot be sorted by date."
	RESOURCE_CANNOT_BE_SORTED_BY_AMOUNT            ExceptionDescription = "Requested resource cannot be sorted by amount."
	RESULTS_NOT_AVAILABLE_FOR_THIS_PAGE            ExceptionDescription = "Requested page results not found."
	NO_SERVICES_RENDERED_IN_THE_PERIOD             ExceptionDescription = "No services were rendered in the requested period."
	INCORRECT_ORDER_STATUS                         ExceptionDescription = "This transaction cannot be performed on this order since it is in '{order_status}' status."
	TRANSACTION_DOES_NOT_EXIST                     ExceptionDescription = "Transaction doesn't exist."
	ORDER_DOES_NOT_EXIST                           ExceptionDescription = "Specified order does not exist."
	COMPANY_ACCOUNT_DOES_NOT_EXIST                 ExceptionDescription = "Specified company account does not exist."
)

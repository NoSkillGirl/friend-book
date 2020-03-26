package constants

const (
	ErrorValidation          = "ValidationError"
	ErrorInternalServerError = "InternalServerError"
	ErrorAlreadyRegistered   = "AlreadyRegistered"
	ErrorStatusUnauthorized  = "Unauthorized"

	ErrorDatabaseConnection             = "DatabaseConnection"
	ErrorDatabaseDuplicate              = "DatabaseDuplicateError"
	ErrorDatabaseInsert                 = "DatabaseInsertError"
	ErrorDatabaseSelect                 = "DatabaseSelectError"
	ErrorDatabaseDelete                 = "DatabaseDeleteError"
	ErrorDatabaseUpdate                 = "DatabaseUpdateError"
	ErrorDatabaseUpdateZeroRowsAffected = "DatabaseUpdateZeroRowsAffected"
	ErrorDatabaseEmailNotFound          = "DatabaseEmailNotFound"
	ErrorDatabaseUserNotFound           = "DatabaseUserNotFound"
)

package errors

import (
	"errors"
	"net/http"
)

type HTTPRequestError struct {
	Err         error
	FTReference string
}

func (e *HTTPRequestError) Error() string {
	return e.Err.Error()
}

func (e *HTTPRequestError) Unwrap() error {
	return e.Err
}

var (
	ErrUnexpected          = errors.New("unexpected error")
	ErrInternalServerError = errors.New("something went wrong. Please try again later.")
	ErrBadRequest          = errors.New("bad request")
	ErrInvalidData         = errors.New("invalid data")
	ErrInvalidAmount       = errors.New("invalid payment amount")
	ErrInvalidCurrency     = errors.New("invalid payment currency")
	ErrInvalidReference    = errors.New("invalid payment reference")
	ErrValidationFailed    = errors.New("request validation failed")

	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidToken = errors.New("invalid or malformed authentication token")
	ErrTokenExpired = errors.New("authentication token has expired")
	ErrForbidden    = errors.New("insufficient permissions")

	ErrPaymentNotFound = errors.New("payment not found")
	ErrUserNotFound    = errors.New("user not found")

	ErrDuplicateReference      = errors.New("payment reference already exists")
	ErrPaymentAlreadyProcessed = errors.New("payment already processed")

	ErrPaymentInitiationFailed   = errors.New("payment initiation failed")
	ErrFundTransferFailed        = errors.New("fund transfer failed")
	ErrTransactionNotFound       = errors.New("transaction not found")
	ErrInvalidPaymentData        = errors.New("invalid payment data")
	ErrTransactionNotCompleted   = errors.New("transaction not completed")
	ErrTransactionReversalFailed = errors.New("failed to revert transaction")
	ErrDuplicateTransaction      = errors.New("duplicate transaction")
	ErrInsufficientBalance       = errors.New("insufficient balance")

	ErrDatabaseConnectionFailed = errors.New("failed to connect to database")
	ErrDatabasePingFailed       = errors.New("failed to ping database")
	ErrDatabaseOperationFailed  = errors.New("database operation failed")
	ErrMessagePublishingFailed  = errors.New("failed to publish payment processing message")
	ErrMessageQueueError        = errors.New("message queue operation failed")
	ErrTimeoutReachingServer    = errors.New("timeout reaching server")
	ErrServerUnreachable        = errors.New("the service is temporarily unavailable")
	ErrRequestTimeout           = errors.New("request timeout")
	ErrTimeout                  = errors.New("request timeout")

	ErrFailedToCreateRequest  = errors.New("failed to create HTTP request")
	ErrFailedToSendRequest    = errors.New("failed to send HTTP request")
	ErrFailedToReadResponse   = errors.New("failed to read response body")
	ErrNon200Response         = errors.New("non-200 response from endpoint")
	ErrInvalidEndpoint        = errors.New("invalid or non-existent endpoint")
	ErrUnexpectedHTTPError    = errors.New("unexpected HTTP client error")
	ErrUnsupportedContentType = errors.New("unsupported content type")

	ErrFailedToMarshalPayload    = errors.New("failed to marshal payload")
	ErrFailedToMarshalPayloadRes = errors.New("failed to marshal response")
	ErrEmptyPayload              = errors.New("empty payload")
	ErrInvalidDataFormat         = errors.New("invalid data format")
	ErrInvalidParameterName      = errors.New("invalid parameter name")
	ErrInvalidID                 = errors.New("invalid ID")

	ErrPaymentNotFoundCompat         = ErrPaymentNotFound
	ErrPaymentAlreadyProcessedCompat = ErrPaymentAlreadyProcessed
	ErrDuplicateReferenceCompat      = ErrDuplicateReference
	ErrInvalidPaymentStatus          = errors.New("invalid payment status transition")
)

var ErrorMap = map[error]int{
	ErrBadRequest:             http.StatusBadRequest,
	ErrInvalidData:            http.StatusBadRequest,
	ErrInvalidAmount:          http.StatusBadRequest,
	ErrInvalidCurrency:        http.StatusBadRequest,
	ErrInvalidReference:       http.StatusBadRequest,
	ErrValidationFailed:       http.StatusBadRequest,
	ErrInvalidPaymentData:     http.StatusBadRequest,
	ErrInvalidParameterName:   http.StatusBadRequest,
	ErrInvalidID:              http.StatusBadRequest,
	ErrEmptyPayload:           http.StatusBadRequest,
	ErrInvalidDataFormat:      http.StatusBadRequest,
	ErrUnsupportedContentType: http.StatusUnsupportedMediaType,

	ErrUnauthorized: http.StatusUnauthorized,
	ErrInvalidToken: http.StatusUnauthorized,
	ErrTokenExpired: http.StatusUnauthorized,

	ErrForbidden: http.StatusForbidden,

	ErrPaymentNotFound:     http.StatusNotFound,
	ErrUserNotFound:        http.StatusNotFound,
	ErrTransactionNotFound: http.StatusNotFound,

	ErrDuplicateReference:      http.StatusConflict,
	ErrPaymentAlreadyProcessed: http.StatusConflict,
	ErrDuplicateTransaction:    http.StatusConflict,

	ErrUnexpected:                http.StatusInternalServerError,
	ErrInternalServerError:       http.StatusInternalServerError,
	ErrPaymentInitiationFailed:   http.StatusInternalServerError,
	ErrFundTransferFailed:        http.StatusInternalServerError,
	ErrTransactionReversalFailed: http.StatusInternalServerError,
	ErrDatabaseConnectionFailed:  http.StatusServiceUnavailable,
	ErrDatabasePingFailed:        http.StatusServiceUnavailable,
	ErrDatabaseOperationFailed:   http.StatusServiceUnavailable,
	ErrMessagePublishingFailed:   http.StatusInternalServerError,
	ErrMessageQueueError:         http.StatusInternalServerError,
	ErrFailedToMarshalPayload:    http.StatusInternalServerError,
	ErrFailedToMarshalPayloadRes: http.StatusInternalServerError,
	ErrFailedToCreateRequest:     http.StatusInternalServerError,
	ErrFailedToSendRequest:       http.StatusBadGateway,
	ErrFailedToReadResponse:      http.StatusInternalServerError,

	ErrRequestTimeout:        http.StatusRequestTimeout,
	ErrTimeout:               http.StatusGatewayTimeout,
	ErrTimeoutReachingServer: http.StatusGatewayTimeout,

	ErrServerUnreachable:   http.StatusServiceUnavailable,
	ErrNon200Response:      http.StatusBadGateway,
	ErrInvalidEndpoint:     http.StatusBadRequest,
	ErrUnexpectedHTTPError: http.StatusBadGateway,

	ErrInsufficientBalance:     http.StatusBadRequest,
	ErrTransactionNotCompleted: http.StatusBadRequest,
}

var UserFriendlyMessages = map[error]string{
	ErrUnexpected:                "An unexpected error occurred. Please try again later.",
	ErrInternalServerError:       "Something went wrong. Please try again later.",
	ErrBadRequest:                "The request you sent is invalid. Please check your input and try again.",
	ErrInvalidData:               "The data you provided is invalid. Please check and try again.",
	ErrInvalidAmount:             "The payment amount must be a valid positive decimal number.",
	ErrInvalidCurrency:           "Currency must be either ETB or USD.",
	ErrInvalidReference:          "Reference must be 1-255 characters long and unique.",
	ErrValidationFailed:          "Request validation failed. Please check your input.",
	ErrUnauthorized:              "You are not authorized to perform this action.",
	ErrInvalidToken:              "Invalid or malformed authentication token.",
	ErrTokenExpired:              "Your authentication token has expired. Please login again.",
	ErrForbidden:                 "You don't have permission to perform this action.",
	ErrPaymentNotFound:           "Payment not found. Please check the payment ID.",
	ErrUserNotFound:              "User not found. Please check your credentials.",
	ErrDuplicateReference:        "A payment with this reference already exists. Please use a different reference.",
	ErrPaymentAlreadyProcessed:   "This payment has already been processed.",
	ErrPaymentInitiationFailed:   "Unable to initiate payment. Please try again.",
	ErrFundTransferFailed:        "Payment transfer failed. Please try again or contact support.",
	ErrTransactionNotFound:       "Transaction not found. Please check your transaction details.",
	ErrInvalidPaymentData:        "Invalid payment information. Please check your details and try again.",
	ErrTransactionNotCompleted:   "Transaction could not be completed. Please try again.",
	ErrTransactionReversalFailed: "Unable to reverse the transaction. Please contact support.",
	ErrDuplicateTransaction:      "Duplicate transaction detected. Please try again.",
	ErrInsufficientBalance:       "Insufficient balance. Please check your account balance.",
	ErrDatabaseConnectionFailed:  "Service temporarily unavailable. Please try again later.",
	ErrDatabasePingFailed:        "Service temporarily unavailable. Please try again later.",
	ErrDatabaseOperationFailed:   "Service temporarily unavailable. Please try again later.",
	ErrMessagePublishingFailed:   "Unable to queue payment for processing. Please try again.",
	ErrMessageQueueError:         "Message processing temporarily unavailable. Please try again later.",
	ErrRequestTimeout:            "Your request is taking longer than expected. Please try again.",
	ErrTimeout:                   "The request timed out. Please try again.",
	ErrTimeoutReachingServer:     "The service is taking longer than expected. Please try again.",
	ErrServerUnreachable:         "The service is temporarily unavailable. Please try again later.",
	ErrFailedToCreateRequest:     "Unable to process your request. Please try again.",
	ErrFailedToSendRequest:       "Unable to send your request. Please try again.",
	ErrFailedToReadResponse:      "Unable to process the response. Please try again.",
	ErrNon200Response:            "The service is temporarily unavailable. Please try again later.",
	ErrInvalidEndpoint:           "Invalid service endpoint. Please contact support.",
	ErrUnexpectedHTTPError:       "We're experiencing technical difficulties. Please try again later.",
	ErrUnsupportedContentType:    "The content type is not supported.",
	ErrFailedToMarshalPayload:    "There was an error processing your request.",
	ErrFailedToMarshalPayloadRes: "There was an error processing the response.",
	ErrEmptyPayload:              "Required information is missing. Please provide all required details.",
	ErrInvalidDataFormat:         "Invalid data format. Please check your input and try again.",
	ErrInvalidParameterName:      "Invalid parameter provided.",
	ErrInvalidID:                 "The ID you provided is invalid.",
}

func GetUserFriendlyMessage(err error) string {
	if err == nil {
		return "An unexpected error occurred. Please try again later."
	}

	if userMsg, exists := UserFriendlyMessages[err]; exists {
		return userMsg
	}

	for targetErr, userMsg := range UserFriendlyMessages {
		if errors.Is(err, targetErr) {
			if err.Error() != targetErr.Error() {
				return err.Error()
			}
			return userMsg
		}
	}

	return err.Error()
}

func GetHTTPStatus(err error) int {
	if status, exists := ErrorMap[err]; exists {
		return status
	}
	return http.StatusInternalServerError
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
	Code   string `json:"code,omitempty"`
}

func CreateErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:  GetUserFriendlyMessage(err),
		Status: GetHTTPStatus(err),
		Code:   GetErrorCode(err),
	}
}

func GetErrorCode(err error) string {
	switch err {
	case ErrValidationFailed, ErrBadRequest, ErrInvalidData, ErrInvalidAmount, ErrInvalidCurrency, ErrInvalidReference, ErrInvalidPaymentData, ErrInvalidParameterName, ErrInvalidID, ErrEmptyPayload, ErrInvalidDataFormat, ErrUnsupportedContentType:
		return "VALIDATION_ERROR"
	case ErrUnauthorized, ErrInvalidToken, ErrTokenExpired:
		return "AUTHENTICATION_ERROR"
	case ErrForbidden:
		return "AUTHORIZATION_ERROR"
	case ErrPaymentNotFound, ErrUserNotFound, ErrTransactionNotFound:
		return "NOT_FOUND_ERROR"
	case ErrDuplicateReference, ErrPaymentAlreadyProcessed, ErrDuplicateTransaction:
		return "CONFLICT_ERROR"
	case ErrPaymentInitiationFailed, ErrFundTransferFailed, ErrTransactionReversalFailed, ErrTransactionNotCompleted:
		return "PAYMENT_ERROR"
	case ErrInsufficientBalance:
		return "BUSINESS_LOGIC_ERROR"
	case ErrDatabaseConnectionFailed, ErrDatabasePingFailed, ErrDatabaseOperationFailed, ErrMessagePublishingFailed, ErrMessageQueueError:
		return "SERVICE_UNAVAILABLE"
	case ErrRequestTimeout, ErrTimeout, ErrTimeoutReachingServer, ErrServerUnreachable, ErrFailedToCreateRequest, ErrFailedToSendRequest, ErrFailedToReadResponse, ErrNon200Response, ErrInvalidEndpoint, ErrUnexpectedHTTPError:
		return "NETWORK_ERROR"
	case ErrUnexpected, ErrInternalServerError, ErrFailedToMarshalPayload, ErrFailedToMarshalPayloadRes:
		return "INTERNAL_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

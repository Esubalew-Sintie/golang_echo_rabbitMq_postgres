package response

import (
	"net/http"
	"payment-gateway/internal/constant/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message,omitempty"`
	Data    any                    `json:"data,omitempty"`
	Error   *DetailedErrorResponse `json:"error,omitempty"`
	Meta    any                    `json:"meta,omitempty"`
}

type DetailedErrorResponse struct {
	Type        string       `json:"type"`
	Message     string       `json:"message"`
	Detail      []FieldError `json:"details,omitempty"`
	FTReference string       `json:"ft_reference,omitempty"`
}

type FieldError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ErrorFields(err error) []FieldError {
	var errs []FieldError

	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			nestedErrors := ErrorFields(v)
			if len(nestedErrors) > 0 {
				errs = append(errs, nestedErrors...)
			} else {
				errs = append(errs, FieldError{
					Name:        i,
					Description: v.Error(),
				})
			}
		}

		return errs
	}

	return nil
}

func getErrorType(err error) string {
	switch err {
	case errors.ErrValidationFailed, errors.ErrInvalidAmount, errors.ErrInvalidCurrency, errors.ErrInvalidReference:
		return "validation_error"
	case errors.ErrUnauthorized, errors.ErrInvalidToken, errors.ErrTokenExpired:
		return "authentication_error"
	case errors.ErrForbidden:
		return "authorization_error"
	case errors.ErrPaymentNotFound, errors.ErrUserNotFound, errors.ErrTransactionNotFound:
		return "not_found_error"
	case errors.ErrDuplicateReference, errors.ErrPaymentAlreadyProcessed, errors.ErrDuplicateTransaction:
		return "conflict_error"
	case errors.ErrPaymentInitiationFailed, errors.ErrFundTransferFailed, errors.ErrTransactionReversalFailed, errors.ErrTransactionNotCompleted:
		return "payment_error"
	case errors.ErrInsufficientBalance:
		return "business_logic_error"
	case errors.ErrDatabaseConnectionFailed, errors.ErrDatabasePingFailed, errors.ErrDatabaseOperationFailed, errors.ErrMessagePublishingFailed, errors.ErrMessageQueueError:
		return "service_unavailable"
	case errors.ErrRequestTimeout, errors.ErrTimeout, errors.ErrTimeoutReachingServer, errors.ErrServerUnreachable, errors.ErrFailedToCreateRequest, errors.ErrFailedToSendRequest, errors.ErrFailedToReadResponse, errors.ErrNon200Response, errors.ErrInvalidEndpoint, errors.ErrUnexpectedHTTPError:
		return "network_error"
	case errors.ErrUnexpected, errors.ErrInternalServerError, errors.ErrFailedToMarshalPayload, errors.ErrFailedToMarshalPayloadRes:
		return "internal_error"
	default:
		return "unknown_error"
	}
}

func SendSuccessResponse(c echo.Context, statusCode int, message string, data any, metadata any) error {
	return c.JSON(statusCode, Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
		Meta:    metadata,
	})
}

func SendErrorResponse(c echo.Context, err error) error {
	if ve, ok := err.(validation.Errors); ok {
		fieldErr := ErrorFields(ve)
		return c.JSON(http.StatusBadRequest, Response{
			Status: http.StatusBadRequest,
			Error: &DetailedErrorResponse{
				Type:    "validation",
				Message: "Please check your input and try again.",
				Detail:  fieldErr,
			},
		})
	}

	userFriendlyMsg := errors.GetUserFriendlyMessage(err)
	statusCode := errors.GetHTTPStatus(err)
	errorType := getErrorType(err)

	ftRef := ""
	if httpErr, ok := err.(*errors.HTTPRequestError); ok {
		ftRef = httpErr.FTReference
	}

	return c.JSON(statusCode, Response{
		Status: statusCode,
		Error: &DetailedErrorResponse{
			Type:        errorType,
			Message:     userFriendlyMsg,
			FTReference: ftRef,
		},
	})
}

func SendValidationErrorResponse(c echo.Context, message string, fieldErrors []FieldError) error {
	return c.JSON(http.StatusBadRequest, Response{
		Status: http.StatusBadRequest,
		Error: &DetailedErrorResponse{
			Type:    "validation_error",
			Message: message,
			Detail:  fieldErrors,
		},
	})
}

func SendValidationError(c echo.Context, validationErr error) error {
	statusCode := errors.GetHTTPStatus(validationErr)
	userFriendlyMsg := errors.GetUserFriendlyMessage(validationErr)
	errorType := getErrorType(validationErr)

	ftRef := ""
	if httpErr, ok := validationErr.(*errors.HTTPRequestError); ok {
		ftRef = httpErr.FTReference
	}

	return c.JSON(statusCode, Response{
		Status: statusCode,
		Error: &DetailedErrorResponse{
			Type:        errorType,
			Message:     userFriendlyMsg,
			FTReference: ftRef,
		},
	})
}

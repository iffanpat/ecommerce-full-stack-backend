package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse represents a standard API response structure
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, code, message, details string) {
	c.JSON(statusCode, StandardResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, details string) {
	ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request data", details)
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *gin.Context, resource string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", resource+" not found", "")
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, details string) {
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", details)
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", message, "")
}

// ConflictResponse sends a conflict error response
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message, "")
}

// UnprocessableEntityResponse sends an unprocessable entity error response
func UnprocessableEntityResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message, "")
}

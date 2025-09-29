package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// --- Success Responses ---

func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, data interface{}, meta interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, APIResponse{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Updated(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Deleted(c *gin.Context, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: message,
	})
}

// --- Empty / No Content ---
func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{})
}

// --- Error Responses ---
func Error(c *gin.Context, code int, message string, err interface{}) {
	c.JSON(code, APIResponse{
		Code:    code,
		Status:  "error",
		Message: message,
		Error:   err,
	})
}

func BadRequest(c *gin.Context, err interface{}) {
	Error(c, http.StatusBadRequest, "Bad request", err)
}

func Unauthorized(c *gin.Context, err interface{}) {
	Error(c, http.StatusUnauthorized, "Unauthorized", err)
}

func Forbidden(c *gin.Context, err interface{}) {
	Error(c, http.StatusForbidden, "Forbidden", err)
}

func NotFound(c *gin.Context, err interface{}) {
	Error(c, http.StatusNotFound, "Resource not found", err)
}

func InternalError(c *gin.Context, err interface{}) {
	Error(c, http.StatusInternalServerError, "Internal server error", err)
}
package utils

import (
	"github.com/gin-gonic/gin"
	"log"
)

func SetError(e interface{}, c *gin.Context, output *Output, internalErrorCode int, externalErrorCode int) {
	var err Error
	var errorMsg string
	switch e.(type) {
	case error:
		errorMsg = e.(error).Error()
		err.ErrorCode = internalErrorCode
		err.ErrorMsg = errorMsg

	case string:
		err.ErrorMsg = e.(string)
		err.ErrorCode = internalErrorCode
	}
	output.Error = append(output.Error, err)

	output.Status = externalErrorCode
	output.Message = "خطا"

	c.JSON(output.Status, output)
	return
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

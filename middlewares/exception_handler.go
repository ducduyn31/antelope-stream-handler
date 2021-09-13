package middlewares

import (
"github.com/gin-gonic/gin"
"net/http"
)
type HttpException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (exception HttpException) Error() string {
	return exception.Message
}

func HttpExceptionHandle() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		detectedErrors := context.Errors.ByType(gin.ErrorTypeAny)

		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err

			var parsedError *HttpException

			switch err.(type) {
			case *HttpException:
				parsedError = err.(*HttpException)
			default:
				parsedError = &HttpException{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				}
			}

			context.AbortWithStatusJSON(parsedError.Code, parsedError)
			return
		}
	}
}

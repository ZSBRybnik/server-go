package routes

import (
	"errors"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

// VerifyTokenHandler - Handling verify-token route
func VerifyTokenHandler(context *gin.Context) {
	token := context.GetHeader("Authorization")
	err := utils.VerifyToken(token)
	utils.ErrorHandler(err, false)
	if err != nil {
		context.AbortWithError(403, errors.New("Forbidden"))
	} else {
		context.JSON(200, "Ok!")
	}
}

package routes

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"zsbrybnik.pl/server-go/utils"
)

type changePasswordJSON struct {
	Token    string `json:"token"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// ChangePasswordHandler - Handling change-password route
func ChangePasswordHandler(context *gin.Context) {
	var changePassword changePasswordJSON
	err := context.Bind(&changePassword)
	utils.ErrorHandler(err, false)
	if err != nil {
		context.AbortWithError(400, errors.New("Bad Request"))
	} else {
		err := utils.VerifyToken(changePassword.Token)
		utils.ErrorHandler(err, false)
		if err != nil {
			context.AbortWithError(403, errors.New("Forbidden"))
		} else {
			passwordInBytes := []byte(changePassword.Password)
			hashedPasswordInBytes, err := bcrypt.GenerateFromPassword(passwordInBytes, 15)
			utils.ErrorHandler(err, false)
			if err != nil {
				context.AbortWithError(500, errors.New("Internal Server Error"))
			} else {
				hashedPassword := string(hashedPasswordInBytes)
				database, ok := context.MustGet("database").(*sql.DB)
				if !ok {
					log.Fatalln("Can't find database in gin-gonic context")
					context.AbortWithError(500, errors.New("Internal Server Error"))
				} else {
					query := "UPDATE admins SET password = ? WHERE login = ?"
					result, err := database.Query(query, hashedPassword, changePassword.Login)
					utils.ErrorHandler(err, false)
					defer result.Close()
					if err != nil {
						context.AbortWithError(500, errors.New("Internal Server Error"))
					} else {
						context.JSON(200, gin.H{
							"status": "Ok!",
						})
					}
				}
			}
		}
	}
}

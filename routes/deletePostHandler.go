package routes

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

type deletePostJSON struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

// DeletePostHandler - Handling delete-post route
func DeletePostHandler(context *gin.Context) {
	var deletePostData deletePostJSON
	err := context.Bind(&deletePostData)
	utils.ErrorHandler(err, false)
	if err != nil {
		context.AbortWithError(400, errors.New("Bad Request"))
	} else {
		err := utils.VerifyToken(deletePostData.Token)
		utils.ErrorHandler(err, false)
		if err != nil {
			context.AbortWithError(403, errors.New("Forbidden"))
		} else {
			database, ok := context.MustGet("database").(*sql.DB)
			if !ok {
				log.Fatalln("Can't find database in gin-gonic context")
				context.AbortWithError(500, errors.New("Internal Server Error"))
			} else {
				query := "DELETE FROM posts WHERE id = ?"
				result, err := database.Query(query, deletePostData.ID)
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

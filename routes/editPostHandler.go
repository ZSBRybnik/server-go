package routes

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

type editPostJSON struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Img          string `json:"img"`
	Introduction string `json:"introduction"`
	ImgAlt       string `json:"imgAlt"`
	Token        string `json:"token"`
	Content      string `json:"content"`
}

// EditPostHandler - Handling edit-post route
func EditPostHandler(context *gin.Context) {
	var editPostData editPostJSON
	err := context.Bind(&editPostData)
	utils.ErrorHandler(err, false)
	if err != nil {
		context.AbortWithError(400, errors.New("Bad Request"))
	} else {
		err := utils.VerifyToken(editPostData.Token)
		utils.ErrorHandler(err, false)
		if err != nil {
			context.AbortWithError(403, errors.New("Forbidden"))
		} else {
			database, ok := context.MustGet("database").(*sql.DB)
			if !ok {
				log.Fatalln("Can't find database in gin-gonic context")
				context.AbortWithError(500, errors.New("Internal Server Error"))
			} else {
				query := "UPDATE posts SET title = ?, content = ?, img = ?, img_alt = ?, introduction = ? WHERE id = ?"
				result, err := database.Query(query, editPostData.Title, editPostData.Content, editPostData.Img, editPostData.ImgAlt, editPostData.Introduction, editPostData.ID)
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

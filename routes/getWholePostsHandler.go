package routes

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

type getWholePostsJSON struct {
	Title        string `json:"title"`
	Content      string `json:"content"`
	Img          string `json:"img"`
	ImgAlt       string `json:"imgAlt"`
	Introduction string `json:"introduction"`
}

// GetWholePostsHandler - Handling get-whole-posts routes
func GetWholePostsHandler(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	database, ok := context.MustGet("database").(*sql.DB)
	if !ok {
		log.Fatalln("Can't find database in gin-gonic context")
		context.AbortWithError(500, errors.New("Internal Server Error"))
	} else {
		query := "SELECT title, content, img, img_alt AS imgAlt, introduction FROM posts WHERE id = ?"
		result := database.QueryRow(query, id)
		var getWholePosts getWholePostsJSON
		err := result.Scan(&getWholePosts.Title, &getWholePosts.Content, &getWholePosts.Img, &getWholePosts.ImgAlt, &getWholePosts.Introduction)
		utils.ErrorHandler(err, false)
		context.JSON(200, gin.H{
			"data": getWholePosts,
		})
	}
}

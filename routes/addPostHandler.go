package routes

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

type addPolishPostJSON struct {
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
	Img          string `json:"img"`
	ImgAlt       string `json:"imgAlt"`
	Content      string `json:"content"`
	Author       string `json:"author"`
}

type addNotPolishPostJSON struct {
	addPolishPostJSON
	Language string `json:"language"`
	PostID   string `json:"postID"`
}

type addPostJSON struct {
	addPolishPostJSON
	addNotPolishPostJSON
}

// AddPostHandler - Handling add-post route
func AddPostHandler(context *gin.Context) {
	action := context.Query("action")
	token := context.GetHeader("Authorization")
	if action == "addPolishPost" || action == "addNotPolishPost" {
		var addPostData addPostJSON
		err := context.Bind(&addPostData)
		utils.ErrorHandler(err, false)
		if err == nil {
			err := utils.VerifyToken(token)
			utils.ErrorHandler(err, false)
			if err == nil {
				database, ok := context.MustGet("database").(*sql.DB)
				if ok {
					var query string
					var result *sql.Rows
					if action == "addPolishPost" {
						query = "INSERT INTO posts (post_polish_id, title, introduction, content, img, img_alt, author, language) SELECT MAX(post_polish_id) + 1 as highestPostId, ?, ?, ?, ?, ?, ?, \"pl\" from posts"
						result, err = database.Query(query, addPostData.Title, addPostData.Introduction, addPostData.Content, addPostData.Img, addPostData.ImgAlt, addPostData.Author)
					} else {
						postIDAsNumber, err := strconv.ParseUint(addPostData.PostID, 10, 32)
						if err == nil {
							query = "INSERT INTO posts (post_polish_id, title, introduction, content, img, img_alt, author, language) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
							result, err = database.Query(query, postIDAsNumber, addPostData.Title, addPostData.Introduction, addPostData.Content, addPostData.Img, addPostData.ImgAlt, addPostData.Author, addPostData.Language)
							utils.ErrorHandler(err, false)
						} else {
							context.AbortWithError(400, errors.New("Bad Request"))
							return
						}
					}
					utils.ErrorHandler(err, false)
					if err != nil {
						context.AbortWithError(500, errors.New("Internal Server Error"))
					} else {
						defer result.Close()
						context.JSON(200, gin.H{
							"status": "Ok!",
						})
						return
					}
				} else {
					log.Fatalln("Can't find database in gin-gonic context")
					context.AbortWithError(500, errors.New("Internal Server Error"))
				}
			} else {
				context.AbortWithError(403, errors.New("Forbidden"))
			}
		} else {
			context.AbortWithError(400, errors.New("Bad Request"))
		}
	} else {
		context.AbortWithError(400, errors.New("Bad Request"))
	}
}

package routes

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

type postJSON struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Img          string `json:"img"`
	Introduction string `json:"introduction"`
	ImgAlt       string `json:"imgAlt"`
}

// GetPostsHandler - Handling get-post route
func GetPostsHandler(context *gin.Context) {
	toSubtract := context.DefaultQuery("toSubtract", "0")
	language := context.Query("language")
	database, ok := context.MustGet("database").(*sql.DB)
	if !ok {
		log.Fatalln("Can't find database in gin-gonic context")
		context.AbortWithError(500, errors.New("Internal Server Error"))
	} else {
		toSubstractAsNumber, err := strconv.Atoi(toSubtract)
		utils.ErrorHandler(err, false)
		if err != nil {
			context.AbortWithError(400, errors.New("Bad Request"))
		} else {
			query := "SELECT post_polish_id, title, img, introduction, img_alt as imgAlt FROM posts WHERE post_polish_id > ((SELECT MAX(post_polish_id) as highestId FROM posts) - ?) AND post_polish_id <= ((SELECT MAX(post_polish_id) as highestId FROM posts) - ?) AND language = \"pl\" ORDER BY post_polish_id DESC"
			toSubstractBorderPost := strconv.Itoa(toSubstractAsNumber + 10)
			result, err := database.Query(query, toSubstractBorderPost, toSubtract)
			utils.ErrorHandler(err, false)
			defer result.Close()
			if err != nil {
				context.AbortWithError(500, errors.New("Internal Server Error"))
			} else {
				var postsArray []postJSON
				for result.Next() {
					var post postJSON
					err := result.Scan(&post.ID, &post.Title, &post.Img, &post.Introduction, &post.ImgAlt)
					utils.ErrorHandler(err, false)
					postsArray = append(postsArray, post)
				}
				if language != "pl" {
					query := "SELECT post_polish_id, title, img, introduction, img_alt as imgAlt FROM posts WHERE post_polish_id > ((SELECT MAX(post_polish_id) as highestId FROM posts) - ?) AND post_polish_id <= ((SELECT MAX(post_polish_id) as highestId FROM posts) - ?) AND language = ? ORDER BY post_polish_id DESC"
					result, err := database.Query(query, toSubstractBorderPost, toSubtract, language)
					utils.ErrorHandler(err, false)
					defer result.Close()
					for result.Next() {
						var post postJSON
						err := result.Scan(&post.ID, &post.Title, &post.Img, &post.Introduction, &post.ImgAlt)
						utils.ErrorHandler(err, false)
						for i, value := range postsArray {
							if value.ID == post.ID {
								postsArray[i].Img = post.Img
								postsArray[i].ImgAlt = post.ImgAlt
								postsArray[i].Introduction = post.Introduction
								postsArray[i].Title = post.Title
							}
						}
						postsArray = append(postsArray, post)
					}
				}
				context.JSON(200, postsArray)
			}
		}
	}
}

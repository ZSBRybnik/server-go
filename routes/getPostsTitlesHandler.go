package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"zsbrybnik.pl/server-go/utils"
)

type postsTitlesJSON struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// GetPostsTitlesHandler Handling get-posts routes
func GetPostsTitlesHandler(context *gin.Context) {
	action := context.Query("action")
	if action == "getPolishPostsTitles" || action == "getNotPolishPostsTitles" {
		redisDB, ok := context.MustGet("redisDB").(*redis.Client)
		if ok {
			var postsTitlesArray []postsTitlesJSON
			value, err := getPostsTitlesCache(context, redisDB, action)
			if err == nil && value != "null" {
				err = json.Unmarshal([]byte(value), &postsTitlesArray)
				utils.ErrorHandler(nil, false)
				context.JSON(200, postsTitlesArray)
			} else {
				database, ok := context.MustGet("database").(*sql.DB)
				if ok {
					query := setupPostsTitlesQuery(action)
					result, err := database.Query(query)
					utils.ErrorHandler(err, false)
					defer result.Close()
					if err == nil {
						postsTitlesArray = scanPostsTitlesResult(result, func(err error) {
							utils.ErrorHandler(err, false)
						})
						redisKey := setupPostsTitlesRedisKey(action)
						postsTitlesArrayInBytes, err := json.Marshal(postsTitlesArray)
						utils.ErrorHandler(err, false)
						redisDB.Set(utils.AppContext, redisKey, postsTitlesArrayInBytes, 10*time.Minute)
						context.JSON(200, postsTitlesArray)
					} else {
						context.AbortWithError(500, errors.New("Internal Server Error"))
					}
				} else {
					log.Fatalln("Can't find database in gin-gonic context")
					context.AbortWithError(500, errors.New("Internal Server Error"))
				}
			}
		} else {
			log.Fatalln("Can't find redis in gin-gonic context")
			context.AbortWithError(500, errors.New("Internal Server Error"))
		}
	} else {
		context.AbortWithError(400, errors.New("Bad Request"))
	}
}

func getPostsTitlesCache(context *gin.Context, redisDB *redis.Client, action string) (string, error) {
	var value string
	var err error
	if action == "getPolishPostsTitles" {
		value, err = redisDB.Get(utils.AppContext, "allPolishPostsTitles").Result()
	} else {
		value, err = redisDB.Get(utils.AppContext, "allNotPolishPostsTitles").Result()
	}
	return value, err
}

func setupPostsTitlesRedisKey(action string) string {
	var redisKey string
	if action == "getPolishPostsTitles" {
		redisKey = "allPolishPostsTitles"
	} else {
		redisKey = "allNotPolishPostsTitles"
	}
	return redisKey
}

func scanPostsTitlesResult(result *sql.Rows, errorHandler func(err error)) []postsTitlesJSON {
	var postsTitlesArray []postsTitlesJSON
	for result.Next() {
		var postsTitles postsTitlesJSON
		err := result.Scan(&postsTitles.ID, &postsTitles.Title)
		errorHandler(err)
		postsTitlesArray = append(postsTitlesArray, postsTitles)
	}
	return postsTitlesArray
}

func setupPostsTitlesQuery(action string) string {
	var query string
	if action == "getPolishPostsTitles" {
		query = "SELECT post_polish_id, title FROM posts WHERE language = \"pl\" ORDER BY post_polish_id DESC"
	} else {
		query = "SELECT post_polish_id, title FROM posts WHERE language <> \"pl\" ORDER BY post_polish_id DESC"
	}
	return query
}

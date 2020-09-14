package routes

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

type getSubpageJSON struct {
	Route string `json:"route"`
}

type subpageDataJSON struct {
	Title        string `json:"title"`
	DisplayTitle string `json:"displayTitle"`
	Content      string `json:"content"`
}

// GetSubpageHandler - Handling get-subpage route
func GetSubpageHandler(context *gin.Context) {
	route := context.Query("route")
	language := context.Query("language")
	database, ok := context.MustGet("database").(*sql.DB)
	if !ok {
		log.Fatalln("Can't find database in gin-gonic context")
		context.AbortWithError(500, errors.New("Internal Server Error"))
	} else {
		query := "SELECT title, display_title AS displayTitle, content FROM subpages WHERE route = ? AND language = ?"
		result := database.QueryRow(query, route, language)
		var subpageData subpageDataJSON
		err := result.Scan(&subpageData.Title, &subpageData.DisplayTitle, &subpageData.Content)
		utils.ErrorHandler(err, false)
		if err != nil && language != "pl" {
			query := "SELECT title, display_title AS displayTitle, content FROM subpages WHERE route = ? AND language = \"pl\""
			result := database.QueryRow(query, route)
			err := result.Scan(&subpageData.Title, &subpageData.DisplayTitle, &subpageData.Content)
			utils.ErrorHandler(err, false)
			if err != nil {
				context.AbortWithError(404, errors.New("Internal Server Error"))
			} else {
				context.JSON(200, subpageData)
			}
		} else if err != nil {
			context.AbortWithError(404, errors.New("Internal Server Error"))
		} else {
			context.JSON(200, subpageData)
		}
	}
}

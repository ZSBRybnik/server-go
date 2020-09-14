package routes

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

type addSubpageJSON struct {
	Token        string `json:"token"`
	Route        string `json:"route"`
	DisplayTitle string `json:"displayTitle"`
}

// AddSubpageHandler - Handling add-subpage routes
func AddSubpageHandler(context *gin.Context) {
	var addSubpageData addSubpageJSON
	err := context.Bind(&addSubpageData)
	utils.ErrorHandler(err, false)
	if err != nil {
		context.AbortWithError(400, errors.New("Bad Request"))
	} else {
		err := utils.VerifyToken(addSubpageData.Token)
		utils.ErrorHandler(err, false)
		if err != nil {
			context.AbortWithError(403, errors.New("Forbidden"))
		} else {
			database, ok := context.MustGet("database").(*sql.DB)
			if !ok {
				log.Fatalln("Can't find database in gin-gonic context")
				context.AbortWithError(500, errors.New("Internal Server Error"))
			} else {
				query := "INSERT INTO subpages (route, display_title) VALUES (?, ?)"
				result, err := database.Query(query, addSubpageData.Route, addSubpageData.DisplayTitle)
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

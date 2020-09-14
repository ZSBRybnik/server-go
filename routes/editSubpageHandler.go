package routes

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

type editSubpageJSON struct {
	Token        string `json:"token"`
	Route        string `json:"route"`
	Title        string `json:"title"`
	DisplayTitle string `json:"displayTitle"`
	Content      string `json:"content"`
}

// EditSubpageHandler - Handling edit-subpage route
func EditSubpageHandler(context *gin.Context) {
	var editSubpageData editSubpageJSON
	err := context.Bind(&editSubpageData)
	utils.ErrorHandler(err, false)
	if err != nil {
		context.AbortWithError(400, errors.New("Bad Request"))
	} else {
		err := utils.VerifyToken(editSubpageData.Token)
		utils.ErrorHandler(err, false)
		if err != nil {
			context.AbortWithError(403, errors.New("Forbidden"))
		} else {
			database, ok := context.MustGet("database").(*sql.DB)
			if !ok {
				log.Fatalln("Can't find database in gin-gonic context")
				context.AbortWithError(500, errors.New("Internal Server Error"))
			} else {
				query := "UPDATE subpages SET route = ?, title = ?, display_title = ?, content = ? WHERE route = ?"
				result, err := database.Query(query, editSubpageData.Route, editSubpageData.Title, editSubpageData.DisplayTitle, editSubpageData.Content, editSubpageData.Route)
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

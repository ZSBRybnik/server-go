package routes

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"zsbrybnik.pl/server-go/utils"
)

type subpagesRoutesCategoriesJSON struct {
	Name          string `json:"name"`
	Title         string `json:"title"`
	OnlyForMobile bool   `json:"onlyForMobile"`
}

// GetSubpagesRoutesCategoriesHandler - handling get-subpages-routes route
func GetSubpagesRoutesCategoriesHandler(context *gin.Context) {
	language := context.Query("language")
	database, ok := context.MustGet("database").(*sql.DB)
	if !ok {
		log.Fatalln("Can't find database in gin-gonic context")
		context.AbortWithError(500, errors.New("Internal Server Error"))
	} else {
		query := "SELECT name, title, only_for_mobile as onlyForMobile FROM routes_categories WHERE language = \"pl\" ORDER BY sort_order"
		result, err := database.Query(query)
		utils.ErrorHandler(err, false)
		defer result.Close()
		var subpagesRoutesArray []subpagesRoutesCategoriesJSON
		for result.Next() {
			var subpageRoute subpagesRoutesCategoriesJSON
			err := result.Scan(&subpageRoute.Name, &subpageRoute.Title, &subpageRoute.OnlyForMobile)
			utils.ErrorHandler(err, false)
			subpagesRoutesArray = append(subpagesRoutesArray, subpageRoute)
		}
		if language != "pl" {
			query := "SELECT name, title, only_for_mobile as onlyForMobile FROM routes_categories WHERE language = ? ORDER BY sort_order"
			result, err := database.Query(query, language)
			utils.ErrorHandler(err, false)
			defer result.Close()
			for result.Next() {
				var subpageRoute subpagesRoutesCategoriesJSON
				err := result.Scan(&subpageRoute.Name, &subpageRoute.Title, &subpageRoute.OnlyForMobile)
				utils.ErrorHandler(err, false)
				for i, value := range subpagesRoutesArray {
					if value.Name == subpageRoute.Name {
						subpagesRoutesArray[i].Title = subpageRoute.Title
					}
				}
			}
		}
		context.JSON(200, subpagesRoutesArray)
	}
}

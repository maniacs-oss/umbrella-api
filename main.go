package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

var (
	db *sql.DB
)

func main() {

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/segments", Auth(false), getSegments)
		v1.GET("/segments/:id", Auth(false), getSegment)
		v1.POST("/segments", Auth(true), addSegment)
		v1.PUT("/segments/:id", Auth(true), editSegment)
		v1.POST("/segments/:id/status", Auth(true), approveSegment)
		v1.DELETE("/segments/:id", Auth(true), deleteSegment)
		v1.GET("/check_items", Auth(false), getCheckItems)
		v1.GET("/check_items/:id", Auth(false), getCheckItem)
		v1.POST("/check_items", Auth(true), addCheckItem)
		v1.PUT("/check_items/:id", Auth(true), editCheckItem)
		v1.DELETE("/check_items/:id", Auth(true), deleteCheckItem)
		v1.POST("/check_items/:id/status", Auth(true), approveCheckItem)
		v1.GET("/account/login_check", Auth(true), loginCheck)
		v1.POST("/account/login", loginEndpoint)
	}
	r.Run(":8080")

}

func newDb(connection string) *sql.DB {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}

	return db
}

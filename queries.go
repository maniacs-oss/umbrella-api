package main

import (
	"database/sql"
	"errors"
	"os"

	"unicode/utf8"

	"github.com/coopernurse/gorp"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func initDb() *gorp.DbMap {
	var err error
	var db *sql.DB
	if os.Getenv("GO_DEV") == "1" {
		db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/umbrella?charset=utf8")
		// db, err = sql.Open("mysql", "apiuser:mEYP4JKdZeeZVbj5@tcp(localhost:1234)/umbrella?charset=utf8")
	} else {
		db, err = sql.Open("mysql", "apiuser:mEYP4JKdZeeZVbj5@tcp(wikiinstance.ciahzg9ywrfk.eu-west-1.rds.amazonaws.com:3306)/umbrella?charset=utf8")
	}
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Segment{}, "segments").SetKeys(true, "Id")
	dbmap.AddTableWithName(CheckItem{}, "check_items").SetKeys(true, "Id")
	return dbmap
}

func checkUser(c *gin.Context, dbmap *gorp.DbMap) (User, error) {

	var user User
	var err error
	token := c.Request.Header.Get("token")
	if utf8.RuneCountInString(token) > 0 {
		err = dbmap.SelectOne(&user, "select id, name, email, password, token from users where token=?", token)
	}
	return user, err
}

func getAllPublishedSegments(c *gin.Context, dbmap *gorp.DbMap) ([]Segment, error) {
	var segments []Segment
	var err error
	_, err = dbmap.Select(&segments, "select id, title, subtitle, body, category from segments WHERE status=:status ORDER BY id ASC", map[string]interface{}{
		"status": "published",
	})
	return segments, err
}

func getSegmentById(c *gin.Context, dbmap *gorp.DbMap, segmentId int64) (Segment, error) {
	var segment Segment
	var err error
	err = dbmap.SelectOne(&segment, "select id, title, subtitle, body, category, status, created_at, author, approved_at, approved_by from segments WHERE id=:segment_id ORDER BY id ASC", map[string]interface{}{
		"status":     "published",
		"segment_id": segmentId,
	})
	if err != nil && err.Error() == "sql: no rows in result set" {
		c.Fail(404, errors.New("The resource could not be found"))
	}
	return segment, err
}

func getAllPublishedCheckItems(c *gin.Context, dbmap *gorp.DbMap) ([]CheckItem, error) {
	var checkItems []CheckItem
	var err error
	_, err = dbmap.Select(&checkItems, "select id, title, text, value, parent, category from check_items WHERE status=:status ORDER BY id ASC", map[string]interface{}{
		"status": "published",
	})
	return checkItems, err
}

func getCheckItemById(c *gin.Context, dbmap *gorp.DbMap, checkItemId int64) (CheckItem, error) {
	var checkItem CheckItem
	var err error
	err = dbmap.SelectOne(&checkItem, "select id, title, text, value, parent, category, status, created_at, author, approved_at, approved_by from check_items WHERE id=:check_item_id ORDER BY id ASC", map[string]interface{}{
		"status":        "published",
		"check_item_id": checkItemId,
	})
	return checkItem, err
}

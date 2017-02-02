package main

import (
	"database/sql"
	"fmt"
	"github.com/securityfirst/umbrella-api/models"
	"github.com/securityfirst/umbrella-api/utils"
	"strings"
	"time"

	"github.com/gosexy/to"

	"github.com/gin-gonic/gin"
)

func (um *Umbrella) getSegments(c *gin.Context) {
	segmentList, err := um.getAllPublishedSegments(c)
	utils.CheckErr(err)
	um.checkErr(c, err)
	um.JSON(c, 200, segmentList)
}

func (um *Umbrella) AddSegment(c *gin.Context) {
	var json models.Segment
	c.Bind(&json)
	fmt.Printf("%+v", json)
	if json.Title == "" || json.Category < 1 {
		c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
		return
	}
	user := c.MustGet("user").(models.User)
	segment := models.Segment{Title: strings.TrimSpace(json.Title), Body: strings.TrimSpace(json.Body), Category: json.Category, Status: "submitted", CreatedAt: time.Now().Unix(), Author: user.Id}
	switch json.DifficultyString {
	case "advanced":
		segment.Difficulty = 2
	case "expert":
		segment.Difficulty = 3
	default:
		segment.Difficulty = 1
	}
	if user.Role == 1 {
		segment.Status = "published"
		segment.ApprovedBy = user.Id
		segment.ApprovedAt = time.Now().Unix()
	}
	err := um.Db.Insert(&segment)
	um.checkErr(c, err)
	c.JSON(200, segment)
}

func (um *Umbrella) EditSegment(c *gin.Context) {
	var json models.Segment
	c.Bind(&json)
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 && (json.Title != "" || json.SubTitle != "" || json.Body != "" || json.Category != 0) {
		segment, err := um.getSegmentById(c, segmentId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
		}
		if json.Title != "" {
			segment.Title = strings.TrimSpace(json.Title)
		}
		if json.Body != "" {
			segment.Body = strings.TrimSpace(json.Body)
		}
		_, err = um.Db.Update(&segment)
		um.checkErr(c, err)
		c.JSON(200, segment)
		return
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func (um *Umbrella) DeleteSegment(c *gin.Context) {
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 {
		segment, err := um.getSegmentById(c, segmentId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
		}
		_, err = um.Db.Delete(&segment)
		um.checkErr(c, err)
		c.JSON(200, gin.H{"response": "Success"})
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}

package main

import (
	"fmt"
	"github.com/ezeql/wiw-scheduler-api/wiw"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type schedulerService struct {
	repository wiw.Repository
}

func main() {

	repository, _ := wiw.NewMySQLRepo("homestead:secret@tcp(localhost:33060)/wheniwork?charset=utf8&parseTime=True&loc=Local")
	router := gin.Default()
	schedulerAPI := schedulerService{repository}

	authMiddleware := router.Group("/", schedulerAPI.Authorization)

	authMiddleware.POST("/shifts/", schedulerAPI.CreateShift)
	authMiddleware.GET("/shifts/", schedulerAPI.ViewShiftsByDate)

	setIDMiddleware := authMiddleware.Group("/", schedulerAPI.ValidateID)
	setIDMiddleware.GET("/users/:id/shifts/", schedulerAPI.ViewShiftsForUser)
	setIDMiddleware.GET("/users/:id/shifts/summarize", schedulerAPI.SummarizeHoursPerWeek)
	setIDMiddleware.GET("/users/:id/viewColleagues/", schedulerAPI.ViewColleagues)
	setIDMiddleware.GET("/users/:id/managers", schedulerAPI.ViewManagers)
	setIDMiddleware.PUT("/shifts/:id", schedulerAPI.UpdateShift)
	setIDMiddleware.GET("/users/:id", schedulerAPI.ViewUser)

	router.Run(":3001")

}

func (ss *schedulerService) Authorization(c *gin.Context) {

}

func (ss *schedulerService) ValidateID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid User Id")
		c.Abort()
		return
	}
	c.Set("id", userID)
	c.Next()
}

func (ss *schedulerService) ViewShiftsForUser(c *gin.Context) {
	shifts, err := ss.repository.ShiftsForUser(c.MustGet("id").(int))
	if err != nil {
		ss.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"shifts": shifts})

}

func (ss *schedulerService) ViewColleagues(c *gin.Context) {
	results, err := ss.repository.ColleaguesForUser(c.MustGet("id").(int))
	if err != nil {
		ss.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"colleages": results})
}

func (ss *schedulerService) SummarizeHoursPerWeek(c *gin.Context) {
	shifts, err := ss.repository.ShiftsForUser(c.MustGet("id").(int))
	if err != nil {
		ss.handleError(c, err)
		return
	}
	weekHours := wiw.SummarizeShifts(shifts)
	c.JSON(http.StatusOK, gin.H{"weekHours": weekHours})
}

func (ss *schedulerService) ViewManagers(c *gin.Context) {
	managers, err := ss.repository.ManagersForUser(c.MustGet("id").(int))
	if err != nil {
		ss.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"managers": managers})
}

func (ss *schedulerService) CreateShift(c *gin.Context) {

}

func (ss *schedulerService) UpdateShift(c *gin.Context) {
}

func (ss *schedulerService) ViewShiftsByDate(c *gin.Context) {
	from := c.DefaultQuery("from", "2000-01-01 00:00:00")
	to := c.DefaultQuery("to", "2100-01-01 00:00:00")
	shifts, err := ss.repository.ShiftsInRange(from, to)
	if err != nil {
		ss.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shifts})

}

func (ss *schedulerService) ViewUser(c *gin.Context) {
	user, err := ss.repository.UserDetails(c.MustGet("id").(int))
	fmt.Println("%v", err)
	if err != nil {
		ss.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})

}

func (ss *schedulerService) handleError(c *gin.Context, err error) {
	if err := err.Error(); err == "record not found" {
		c.JSON(404, gin.H{"error": "record not found"})
	}

}

package main

import (
	"github.com/ezeql/wiw-scheduler-api/wiw"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type schedulerService struct {
	repository wiw.Repository
}

func main() {

	//move this into an env
	dsn := "homestead:secret@tcp(localhost:33060)/wheniwork?charset=utf8&parseTime=True&loc=Local"

	repository, err := wiw.NewMySQLRepo(dsn)
	if err != nil {
		log.Println("ERROR: Cannot build repository")
		log.Println(e)
		log.Fatal("Exiting...")
		return
	}
	router := gin.Default()
	schedulerAPI := schedulerService{repository}

	authMiddleware := router.Group("/", schedulerAPI.Authorization)
	setIDMiddleware := authMiddleware.Group("/", schedulerAPI.ValidateID)

	authMiddleware.POST("/shifts/", schedulerAPI.CreateOrUpdateShift)
	setIDMiddleware.PUT("/shifts/:id", schedulerAPI.CreateOrUpdateShift)
	authMiddleware.GET("/shifts/", schedulerAPI.ViewShiftsByDate)
	setIDMiddleware.GET("/users/:id/shifts/", schedulerAPI.ViewShiftsForUser)
	setIDMiddleware.GET("/users/:id/colleagues/", schedulerAPI.ViewColleagues)
	setIDMiddleware.GET("/users/:id/managers", schedulerAPI.ViewManagers)
	setIDMiddleware.GET("/users/:id", schedulerAPI.ViewUser)

	router.Run(":3001")
}

func (ss *schedulerService) Authorization(c *gin.Context) {
}

func (ss *schedulerService) ValidateID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,  gin.H{"error": "Invalid User Id"} )
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
	result := gin.H{"shifts": shifts}
	summarize := c.Query("summarize")
	if summarize != "" {
		result["summary"] = wiw.SummarizeShifts(shifts)
	}
	c.JSON(http.StatusOK, result)
}

func (ss *schedulerService) ViewColleagues(c *gin.Context) {
	colleages, err := ss.repository.ColleaguesForUser(c.MustGet("id").(int))
	if err != nil {
		ss.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"colleages": colleages})
}

func (ss *schedulerService) ViewManagers(c *gin.Context) {
	managers, err := ss.repository.ManagersForUser(c.MustGet("id").(int))
	if err != nil {
		ss.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"managers": managers})
}

func (ss *schedulerService) CreateOrUpdateShift(c *gin.Context) {
	shift := &wiw.Shift{}

	if err := c.BindJSON(shift); err != nil {
		ss.handleError(c, err)
		return
	}
	if shift.ManagerID == 0 {
		//TODO: should use current user id
	}

	if c.Request.Method == "POST" {
		if err := ss.repository.CreateShift(shift); err != nil {
			ss.handleError(c, err)
			return
		}
	} else {
		shift.ID = uint(c.MustGet("id").(int))
		if err := ss.repository.UpdateOrCreateShift(shift); err != nil {
			ss.handleError(c, err)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"shift": shift})
}

func (ss *schedulerService) ViewShiftsByDate(c *gin.Context) {
	from := c.DefaultQuery("from", "2000-01-01 00:00:00")
	to := c.DefaultQuery("to", "2100-01-01 00:00:00")
	shifts, err := ss.repository.ShiftsInRange(from, to)
	if err != nil {
		ss.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"shifts": shifts})
}

func (ss *schedulerService) ViewUser(c *gin.Context) {
	user, err := ss.repository.UserDetails(c.MustGet("id").(int))
	if err != nil {
		ss.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (ss *schedulerService) handleError(c *gin.Context, err error) {
		//TODO: return appropiate http code based on error type
		c.JSON(http.StatusNotFound, gin.H{"error": err})
	}
}

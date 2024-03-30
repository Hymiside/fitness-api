package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Hymiside/fitness-api/pkg/models"
	"github.com/Hymiside/fitness-api/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s]	REQUEST: %s %s    STATUS-CODE: %d    LATENSY: %s\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/admin/sign-in", h.signInAdmin)
		auth.POST("/trainer/sign-in", h.signInTrainer)
	}

	api := router.Group("/fitness", h.userIdentity)
	{
		api.POST("/trainer/create", h.createTrainer)
		api.GET("/trainer/:id", nil)
		api.GET("/trainer/list", h.getTrainers)

		api.POST("/client/create", h.createClient)
		api.GET("/client/:id", nil)
		api.GET("/client/list", h.getClients)

		api.POST("/workout/create", h.createWorkout)
		api.GET("/workout/:id", nil)
		api.GET("/workout/list-by-date", h.getWorkoutsByDate)  // ?date=2020-01-01
		api.GET("/workout/list-by-interval", nil)  // ?from=2020-01-01&to=2020-01-02

		api.POST("/workout/type/create", h.createWorkoutType)
		api.GET("/workout/type/list", h.getWorkoutTypes)
	}

	return router
}

func (h *Handler) signInAdmin(c *gin.Context) {
	var data models.Admin

	if err := c.BindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.services.GenerateTokenForAdmin(c, data.Login, data.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPwd) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, dbr.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) signInTrainer(c *gin.Context) {
	var data models.Trainer

	if err := c.BindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.services.GenerateTokenForTrainer(c, data.Token)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) createTrainer(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if data["role"] != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
		return
	}

	var trainerData models.Trainer
	if err := c.BindJSON(&trainerData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.services.CreateTrainer(c, trainerData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) getTrainers(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if data["role"] != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
		return
	}

	trainers, err := h.services.GetTrainers(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, trainers)
}

func (h *Handler) createClient(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if data["role"] != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
		return
	}

	var clientData models.Client
	if err := c.BindJSON(&clientData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.CreateClient(c, clientData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) getClients(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if data["role"] != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
		return
	}

	clients, err := h.services.GetClients(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.AbortWithStatusJSON(http.StatusOK, clients)
}

func (h *Handler) createWorkout(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if data["role"] != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
		return
	}

	var workoutData models.WorkoutRequest
	if err := c.BindJSON(&workoutData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.CreateWorkout(c, workoutData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) createWorkoutType(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if data["role"] != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
		return
	}

	var workoutTypeData models.WorkoutType
	if err := c.BindJSON(&workoutTypeData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.CreateWorkoutType(c, workoutTypeData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) getWorkoutTypes(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if data["role"] != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
		return
	}

	workoutTypes, err := h.services.GetWorkoutTypes(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.AbortWithStatusJSON(http.StatusOK, workoutTypes)
}

func (h *Handler) getWorkoutsByDate(c *gin.Context) {
	dateData, ok := c.GetQuery("date")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "date not found"})
		return
	}
	
	date, _ := time.Parse("2006-01-02T15:04:05Z", dateData)
	t, err := time.Parse("2006-01-02", date.Format("2006-01-02"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workouts, err := h.services.GetWorkoutsByDate(c, t)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, workouts)
}
package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
		api.GET("/admin", h.getAdminByID)  // ?id=1
		api.GET("/admin/list", h.getAdmins)
		api.GET("/admin/delete", h.deleteAdmin) // ?id=1
		api.POST("/admin/create", h.createAdmin)
		api.GET("/admin/type", h.getAdminType)

		api.GET("/trainer", h.getTrainerByID)  // ?id=1
		api.POST("/trainer/create", h.createTrainer)
		api.GET("/trainer/delete", h.deleteTrainer) // ?id=1
		api.GET("/trainer/list", h.getTrainers)
		api.GET("/trainer/cash/day", h.GetCashByDay)
		api.GET("/trainer/cash/month", h.GetCashByMonth)
		

		api.POST("/client/create", h.createClient)
		api.GET("/client", h.getClientByID)  // ?id=1
		api.POST("/client/edit", h.updateClient)
		api.GET("/client/list", h.getClients)

		api.GET("/workout", h.getWorkoutByID)  // ?id=1
		api.POST("/workout/create", h.createWorkout)
		api.POST("/workout/edit", h.updateWorkout)
		api.GET("/workout/delete", h.deleteWorkout) // ?id=1

		api.GET("/workout/change-status", h.changeStatusWorkout)  // ?id=1&status=done
		api.GET("/workout/list-by-date", h.getWorkoutsByDate)  // ?date=2023-12-23T15:04:05Z
		api.GET("/workout/list-by-interval", h.getWorkoutsByInterval)  // ?from=2023-12-23T15:04:05Z&to=2023-12-23T15:04:05Z
		api.GET("/workout/list", h.getWorkouts)  // ?trainer_id=1 or ?client_id=1

		api.GET("/workout/type", h.getWorkoutTypeByID)  // ?id=1
		api.POST("/workout/type/create", h.createWorkoutType)
		api.GET("/workout/type/edit", h.updateWorkoutType)
		api.GET("/workout/type/delete", h.deleteWorkoutType)  // ?id=1
		api.GET("/workout/type/list", h.getWorkoutTypes)
	}

	return router
}

func (h *Handler) getWorkoutTypeByID(c *gin.Context) {
	workoutTypeID, _ := strconv.Atoi(c.Query("id"))

	workoutType, err := h.services.GetWorkoutTypeByID(c, workoutTypeID)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, workoutType)
}

func (h *Handler) getAdminByID(c *gin.Context) {
	adminID, _ := strconv.Atoi(c.Query("id"))

	admin, err := h.services.GetAdminByID(c, adminID)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, admin)
}

func (h *Handler) createAdmin(c *gin.Context) {
	admin := models.Admin{}
	if err := c.BindJSON(&admin); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := h.services.CreateAdmin(c, admin); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) getAdmins(c *gin.Context) {
	admins, err := h.services.GetAdmins(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.AbortWithStatusJSON(http.StatusOK, admins)
}

func (h *Handler) deleteAdmin(c *gin.Context) {
	adminID, _ := strconv.Atoi(c.Query("id"))

	if err := h.services.DeleteAdmin(c, adminID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) GetCashByMonth(c *gin.Context) {
	data, _ := getData(c)
	trainerID := data["userID"].(int)

	cash, err := h.services.GetCashByMonth(c, trainerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"cash": cash})
}

func (h *Handler) GetCashByDay(c *gin.Context) {
	data, _ := getData(c)
	trainerID := data["userID"].(int)

	cash, err := h.services.GetCashByDay(c, trainerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"cash": cash})
}

func (h *Handler) getAdminType(c *gin.Context) {
	data, _ := getData(c)
	adminID := data["userID"].(int)

	adminType, err := h.services.GetAdminType(c, adminID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"type": adminType})
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

func (h *Handler) deleteTrainer(c *gin.Context) {
	trainerIDdata, _ := c.GetQuery("id")
	trainerID, err := strconv.Atoi(trainerIDdata)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.DeleteTrainer(c, trainerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) getTrainers(c *gin.Context) {
	trainers, err := h.services.GetTrainers(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, trainers)
}

func (h *Handler) getTrainerByID(c *gin.Context) {
	data, _ := c.GetQuery("id")
	trainerID, err := strconv.Atoi(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trainer, err := h.services.GetTrainerByID(c, trainerID)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, trainer)
}

func (h *Handler) createClient(c *gin.Context) {
	var clientData models.Client
	if err := c.BindJSON(&clientData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.CreateClient(c, clientData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) updateClient(c *gin.Context) {
	var clientData models.Client
	if err := c.BindJSON(&clientData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.UpdateClient(c, clientData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) getClients(c *gin.Context) {
	clients, err := h.services.GetClients(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.AbortWithStatusJSON(http.StatusOK, clients)
}

func (h *Handler) getClientByID(c *gin.Context) {
	queryData, _ := c.GetQuery("id")
	clientID, err := strconv.Atoi(queryData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.services.GetClientByID(c, clientID)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, client)
}

func (h *Handler) createWorkout(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var workoutData models.WorkoutRequest
	if err := c.BindJSON(&workoutData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workoutData.AdminID = data["userID"].(int)
	err = h.services.CreateWorkout(c, workoutData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) updateWorkout(c *gin.Context) {
	var workoutData models.WorkoutRequest
	if err := c.BindJSON(&workoutData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.UpdateWorkout(c, workoutData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) deleteWorkout(c *gin.Context) {
	queryData, _ := c.GetQuery("id")
	workoutID, err := strconv.Atoi(queryData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.DeleteWorkout(c, workoutID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) getWorkouts(c *gin.Context) {
	var (
		trainerID int
		err error
	)

	queryTrainerID, ok := c.GetQuery("trainer_id")
	if ok {
		trainerID, err = strconv.Atoi(queryTrainerID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	
	var clientID int
	queryClientID, ok := c.GetQuery("client_id")
	if ok {
		clientID, err = strconv.Atoi(queryClientID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	workouts, err := h.services.GetWorkouts(c, trainerID, clientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, workouts)
}

func (h *Handler) createWorkoutType(c *gin.Context) {
	var workoutTypeData models.WorkoutType
	if err := c.BindJSON(&workoutTypeData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.CreateWorkoutType(c, workoutTypeData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) updateWorkoutType(c *gin.Context) {
	var workoutTypeData models.WorkoutType
	if err := c.BindJSON(&workoutTypeData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.UpdateWorkoutType(c, workoutTypeData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) deleteWorkoutType(c *gin.Context) {
	queryData, _ := c.GetQuery("id")
	workoutTypeID, err := strconv.Atoi(queryData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.DeleteWorkoutType(c, workoutTypeID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) getWorkoutTypes(c *gin.Context) {
	workoutTypes, err := h.services.GetWorkoutTypes(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.AbortWithStatusJSON(http.StatusOK, workoutTypes)
}

func (h *Handler) getWorkoutsByDate(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var trainerID int
	if data["role"] != "admin" {
		trainerID = data["userID"].(int)
	}


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

	workouts, err := h.services.GetWorkoutsByDate(c, t, trainerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, workouts)
}

func (h *Handler) getWorkoutsByInterval(c *gin.Context) {
	data, err := getData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var trainerID int
	if data["role"] != "admin" {
		trainerID = data["userID"].(int)
	}

	dateFromData, ok := c.GetQuery("from")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "from not found"})
		return
	}

	dateToData, ok := c.GetQuery("to")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "to not found"})
		return
	}

	dateFrom, _ := time.Parse("2006-01-02T15:04:05Z", dateFromData)
	dateTo, _ := time.Parse("2006-01-02T15:04:05Z", dateToData)

	tFrom, err := time.Parse("2006-01-02", dateFrom.Format("2006-01-02"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tTo, err := time.Parse("2006-01-02", dateTo.Format("2006-01-02"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workouts, err := h.services.GetWorkoutsByInterval(c, tFrom, tTo, trainerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, workouts)
}

func (h *Handler) getWorkoutByID(c *gin.Context) {
	data, _ := c.GetQuery("id")
	workoutID, err := strconv.Atoi(data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workout, err := h.services.GetWorkoutByID(c, workoutID)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, workout)
}

func (h *Handler) changeStatusWorkout(c *gin.Context) {
	workoutIDdata, _ := c.GetQuery("id")
	status, _ := c.GetQuery("status")

	workoutID, err := strconv.Atoi(workoutIDdata)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.services.ChangeStatusWorkout(c, workoutID, status); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
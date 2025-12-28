package handlers

import (
	"events/core/models"
	"events/core/ports"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	defaultUserKafkaTopic    = "users-topic"
	defaultPaymentKafkaTopic = "payment-topic"
	defaultMovieKafkaTopic   = "movie-topic"
)

var (
	emptyKey = []byte{}
)

type EventsHandler struct {
	producer ports.EventProducer
}

func NewEventsHandler(producer ports.EventProducer) *EventsHandler {
	return &EventsHandler{
		producer: producer,
	}
}

func (h *EventsHandler) RegisterRoutes(eventsRouter *gin.RouterGroup) {
	eventsRouter.POST("/user", h.CreateUser)
	eventsRouter.POST("/movie", h.CreateMovie)
	eventsRouter.POST("/payment", h.CreatePayment)
}

func (h *EventsHandler) CreateUser(c *gin.Context) {
	var user models.UserEvent

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEvent := models.NewEvent(models.UserCreated, user)
	encoded, err := userEvent.ToJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.producer.SendEvent(c.Request.Context(), defaultUserKafkaTopic, emptyKey, encoded); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Event produced: %s", models.UserCreated)

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (h *EventsHandler) CreateMovie(c *gin.Context) {
	var movie models.MovieEvent

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movieEvent := models.NewEvent(models.MovieCreated, movie)
	encoded, err := movieEvent.ToJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.producer.SendEvent(c.Request.Context(), defaultMovieKafkaTopic, emptyKey, encoded); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Event produced: %s", models.MovieCreated)

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (h *EventsHandler) CreatePayment(c *gin.Context) {
	var payment models.PaymentEvent

	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentEvent := models.NewEvent(models.PaymentCreated, payment)
	encoded, err := paymentEvent.ToJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.producer.SendEvent(c.Request.Context(), defaultPaymentKafkaTopic, emptyKey, encoded); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Event produced: %s", models.PaymentCreated)

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

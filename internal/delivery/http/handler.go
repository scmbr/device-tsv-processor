package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1_handlers "github.com/scmbr/device-tsv-processor/internal/delivery/http/handlers/v1"
	"github.com/scmbr/device-tsv-processor/internal/usecase"
)

type Handler struct {
	ucs *usecase.UseCases
}

func NewHandler(ucs *usecase.UseCases) *Handler {
	return &Handler{ucs: ucs}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery(), gin.Logger())
	h.initAPI(router)
	return router
}
func (h *Handler) initAPI(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"msg": "pong",
		})
	})
	deviceMessagesHandler := v1_handlers.NewDeviceMessageHandler(h.ucs.DeviceMessage.GetDeviceMessages)
	messages := v1.Group("/messages")
	{
		messages.GET("/:guid", deviceMessagesHandler.ListByGUID)
	}
}

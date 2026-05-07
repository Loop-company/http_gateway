package handler

import (
	"net/http"
	"strconv"

	"github.com/Loop-company/http_gateway/internal/analyticsclient"
	gatewaypb "github.com/Loop-company/http_gateway/proto"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	client *analyticsclient.Client
}

func NewAnalyticsHandler(client *analyticsclient.Client) *AnalyticsHandler {
	return &AnalyticsHandler{client: client}
}

func (h *AnalyticsHandler) SearchEvents(c *gin.Context) {
	limit := parseInt32(c.Query("limit"), 100)
	offset := parseInt32(c.Query("offset"), 0)

	resp, err := h.client.SearchEvents(c.Request.Context(), &gatewaypb.SearchEventsRequest{
		UserId:        c.Query("user_id"),
		EventType:     c.Query("event_type"),
		SourceService: c.Query("source_service"),
		From:          c.Query("from"),
		To:            c.Query("to"),
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AnalyticsHandler) RegistrationsReport(c *gin.Context) {
	resp, err := h.client.RegistrationsReport(c.Request.Context(), c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AnalyticsHandler) LoginReport(c *gin.Context) {
	resp, err := h.client.LoginReport(c.Request.Context(), c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AnalyticsHandler) TopUsersReport(c *gin.Context) {
	limit := parseInt32(c.Query("limit"), 10)

	resp, err := h.client.TopUsersReport(c.Request.Context(), c.Query("from"), c.Query("to"), limit)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func parseInt32(value string, fallback int32) int32 {
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return fallback
	}

	return int32(parsed)
}

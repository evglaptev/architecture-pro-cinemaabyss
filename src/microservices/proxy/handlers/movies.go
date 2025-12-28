package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	monolithURL     string
	microserviceURL string

	migrationPercent int

	httpClient *http.Client
}

func NewProxyHandler(monolithURL, microserviceURL string, migrationPercent int) *ProxyHandler {
	return &ProxyHandler{
		monolithURL:      monolithURL,
		microserviceURL:  microserviceURL,
		migrationPercent: migrationPercent,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func relayToMicroservice(migrationPercent int) bool {
	return rand.Intn(100) < migrationPercent
}

func (h *ProxyHandler) selectBaseURL() string {
	if relayToMicroservice(h.migrationPercent) {
		return h.microserviceURL
	}
	return h.monolithURL
}

func (h *ProxyHandler) proxyRequest(c *gin.Context, targetBaseURL string) {
	pathWithQuery := c.Request.URL.Path + "?" + c.Request.URL.Query().Encode()

	fmt.Println("Relay path: ", pathWithQuery, "Relay target: ", targetBaseURL)

	// Clone original HTTP-request.
	req, err := http.NewRequestWithContext(
		c.Request.Context(),
		c.Request.Method,
		targetBaseURL+pathWithQuery,
		c.Request.Body,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create proxy-request",
		})
		return
	}
	req.Header = c.Request.Header.Clone()

	// Perform HTTP-request to target.
	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "Service unavailable",
			"message": "Unsuccessful request to target service",
		})
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func (h *ProxyHandler) RegisterRoutes(router *gin.RouterGroup) {
	movies := router.Group("/movies")
	{
		movies.Any("", h.ProxyByPartialMigration)
		movies.Any("/*path", h.ProxyByPartialMigration)
	}

	users := router.Group("/users")
	{
		users.Any("", h.Proxy)
		users.Any("/*path", h.Proxy)
	}

	payment := router.Group("/payments")
	{
		payment.Any("", h.Proxy)
		payment.Any("/*path", h.Proxy)
	}

	subscribtion := router.Group("/subscriptions")
	{
		subscribtion.Any("", h.Proxy)
		subscribtion.Any("/*path", h.Proxy)
	}
}

func (h *ProxyHandler) ProxyByPartialMigration(c *gin.Context) {
	h.proxyRequest(c, h.selectBaseURL())
}

func (h *ProxyHandler) Proxy(c *gin.Context) {
	h.proxyRequest(c, h.monolithURL)
}

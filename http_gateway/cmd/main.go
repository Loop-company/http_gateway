package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Loop-company/http_gateway/internal/analyticsclient"
	"github.com/Loop-company/http_gateway/internal/authclient"
	"github.com/Loop-company/http_gateway/internal/handler"
	"github.com/Loop-company/http_gateway/internal/middleware"
	"github.com/Loop-company/http_gateway/internal/userclient"
	"github.com/gin-gonic/gin"
)

func main() {
	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authAddr == "" {
		authAddr = "localhost:50051"
	}

	userAddr := os.Getenv("USER_SERVICE_ADDR")
	if userAddr == "" {
		userAddr = "localhost:50052"
	}

	analyticsAddr := os.Getenv("ANALYTICS_SERVICE_ADDR")
	if analyticsAddr == "" {
		analyticsAddr = "localhost:50053"
	}

	authClient, err := authclient.New(authAddr)
	if err != nil {
		log.Fatalf("failed to create auth client: %v", err)
	}
	defer authClient.Close()

	userClient, err := userclient.New(userAddr)
	if err != nil {
		log.Fatalf("failed to create user client: %v", err)
	}
	defer userClient.Close()

	analyticsClient, err := analyticsclient.New(analyticsAddr)
	if err != nil {
		log.Fatalf("failed to create analytics client: %v", err)
	}
	defer analyticsClient.Close()

	authHandler := handler.NewAuthHandler(authClient)
	userHandler := handler.NewUserHandler(userClient)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsClient)

	r := gin.Default()

	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify", authHandler.Verify)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(authClient))
	{
		protected.POST("/auth/logout", authHandler.Logout)

		protected.GET("/users/:id", userHandler.GetProfile)
		protected.PUT("/users/name", userHandler.UpdateName)

		protected.GET("/analytics/events", analyticsHandler.SearchEvents)
		protected.GET("/analytics/reports/registrations", analyticsHandler.RegistrationsReport)
		protected.GET("/analytics/reports/logins", analyticsHandler.LoginReport)
		protected.GET("/analytics/reports/top-users", analyticsHandler.TopUsersReport)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Gateway starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gateway...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Gateway stopped")
}

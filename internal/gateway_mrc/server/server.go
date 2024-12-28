package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	db "job_search_platform/internal/gateway_mrc/db/sqlc"
	"job_search_platform/internal/gateway_mrc/handlers"
	"job_search_platform/internal/gateway_mrc/middleware"
	"job_search_platform/internal/gateway_mrc/usecases"
	"job_search_platform/pkg/config"
	"job_search_platform/pkg/jwt_token"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	config     config.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker jwt_token.Maker
	httpServer *http.Server
	logger     zerolog.Logger
}

func NewServer(config config.Config, store db.Store, logger zerolog.Logger) (*Server, error) {
	tokenMaker, err := jwt_token.NewJWTMaker(
		config.TokenSymmetricKey,
		config.AccessTokenExpiresIn,
		config.RefreshTokenExpiresIn,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		logger:     logger,
	}
	server.setupRouter()

	server.httpServer = &http.Server{
		Addr:           config.HTTPServerAddress,
		Handler:        server.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{server.config.Origin}, // Укажите домен вашего клиента
		AllowCredentials: true,                           // Разрешить использование учетных данных (например, куки)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
	}

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.New(corsConfig))

	// SESSION
	sessionMiddleware := middleware.SessionMiddleware(server.store)
	router.Use(sessionMiddleware)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Route %s not found", ctx.Request.URL)})
	})
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	//
	usecase := usecases.NewSessionsUsecase(server.store)
	handler := handlers.NewProxyHandler(server.tokenMaker, usecase, server.config)
	server.setupAuthRoutes(router, handler)
	server.setupUsersRoutes(router, handler)
	server.router = router
}

func (server *Server) setupAuthRoutes(router *gin.Engine, handler handlers.ProxyHandler) {
	authMiddleware := middleware.AuthMiddleware(server.tokenMaker, server.store, server.config)
	public := router.Group("/api/v1/auth/public")
	{
		public.POST("/sign-up", func(ctx *gin.Context) {
			handler.ProxyCommonReq(ctx, server.config.UsersMrcUrl)
		})
		public.POST("/sign-in", func(ctx *gin.Context) {
			handler.ProxySignInReq(ctx, server.config.UsersMrcUrl)
		})
	}

	private := router.Group("/api/v1/auth/private")
	private.Use(authMiddleware)
	{
		private.GET("/logout", func(ctx *gin.Context) {
			handler.ProxyLogoutReq(ctx)
		})
		private.POST("/*path", func(ctx *gin.Context) {
			handler.ProxyCommonReq(ctx, server.config.UsersMrcUrl)
		})
		private.DELETE("/*path", func(ctx *gin.Context) {
			handler.ProxyCommonReq(ctx, server.config.UsersMrcUrl)
		})
		private.PUT("/*path", func(ctx *gin.Context) {
			handler.ProxyCommonReq(ctx, server.config.UsersMrcUrl)
		})
	}
}

func (server *Server) setupUsersRoutes(router *gin.Engine, handler handlers.ProxyHandler) {
	authMiddleware := middleware.AuthMiddleware(server.tokenMaker, server.store, server.config)
	private := router.Group("/api/v1/users/private")
	address := server.config.UsersMrcUrl
	private.Use(authMiddleware)
	{
		registerRoutes(private, handler, address)
	}
}

func registerRoutes(group *gin.RouterGroup, handler handlers.ProxyHandler, address string) {
	group.GET("/*path", func(ctx *gin.Context) {
		handler.ProxyCommonReq(ctx, address)
	})
	group.POST("/*path", func(ctx *gin.Context) {
		handler.ProxyCommonReq(ctx, address)
	})
	group.DELETE("/*path", func(ctx *gin.Context) {
		handler.ProxyCommonReq(ctx, address)
	})
	group.PUT("/*path", func(ctx *gin.Context) {
		handler.ProxyCommonReq(ctx, address)
	})
}

func (server *Server) Start() error {
	go func() {
		log.Printf("Starting server on %s\n", server.httpServer.Addr)
		if err := server.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Ожидание сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	// Контекст с таймаутом для завершения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Println("Server exited gracefully")
	return nil
}

func RunGinServer(config config.Config, store db.Store, logger zerolog.Logger) error {
	server, err := NewServer(config, store, logger)
	if err != nil {
		return err
	}

	err = server.Start()
	if err != nil {
		return err
	}
	return nil
}

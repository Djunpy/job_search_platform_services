package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	db "job_search_platform/internal/users_mrc/db/sqlc"
	"job_search_platform/internal/users_mrc/handlers"
	"job_search_platform/internal/users_mrc/routes"
	"job_search_platform/internal/users_mrc/usecases"
	"job_search_platform/pkg/config"
	"job_search_platform/pkg/jwt_token"
	"job_search_platform/pkg/middleware"
	"job_search_platform/pkg/scheduler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	config      config.Config
	store       db.Store
	router      *gin.Engine
	tokenMaker  jwt_token.Maker
	distributor scheduler.TaskDistributor
	httpServer  *http.Server
	logger      zerolog.Logger
}

func NewServer(config config.Config, store db.Store, distributor scheduler.TaskDistributor, logger zerolog.Logger) (*Server, error) {
	tokenMaker, err := jwt_token.NewJWTMaker(
		config.TokenSymmetricKey,
		config.AccessTokenExpiresIn,
		config.RefreshTokenExpiresIn,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:      config,
		store:       store,
		tokenMaker:  tokenMaker,
		distributor: distributor,
		logger:      logger,
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
	//router.Use(middleware.OpenCORSMiddleware())
	//router.Use(middleware.HandleSessionMiddleware(server.store))
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Route %s not found", ctx.Request.URL)})
	})
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := router.Group("api")
	v1 := api.Group("/v1")
	//adminGroup := v1.Group("/admin")
	server.setupAuthRoutes(v1)
	server.setupUsersRoutes(v1)
	server.router = router
}

func (server *Server) setupAuthRoutes(rg *gin.RouterGroup) {
	jwtDeserializer := middleware.JWTDeserializer(server.tokenMaker)
	usecase := usecases.NewAuthUsecase(server.store, server.tokenMaker)
	handler := handlers.NewAuthHandler(usecase, server.distributor)
	route := routes.NewAuthRouter(handler)
	router := rg.Group("/auth")
	public := router.Group("/public")
	private := router.Group("/private")
	private.Use(jwtDeserializer)
	route.InitAuthRouter(public, private)
}

func (server *Server) setupUsersRoutes(rg *gin.RouterGroup) {
	jwtDeserializer := middleware.JWTDeserializer(server.tokenMaker)
	usecase := usecases.NewUsersUsecase(server.store)
	handler := handlers.NewUsersHandler(usecase)
	route := routes.NewUsersRouter(handler)
	router := rg.Group("/users")
	public := router.Group("/public")
	private := router.Group("/private")
	private.Use(jwtDeserializer)
	route.InitUsersRouter(public, private)
}

func (server *Server) Start() error {
	go func() {
		server.logger.Info().Msg(fmt.Sprintf("Starting server on %s\n", server.httpServer.Addr))
		if err := server.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Ожидание сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	server.logger.Info().Msg("Shutting down server...")

	// Контекст с таймаутом для завершения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	server.logger.Info().Msg("Server exited gracefully")
	return nil
}

func RunGinServer(config config.Config, store db.Store, distributor scheduler.TaskDistributor, logger zerolog.Logger) error {
	server, err := NewServer(config, store, distributor, logger)
	if err != nil {
		return err
	}
	err = server.Start()
	if err != nil {
		return err
	}
	return nil
}

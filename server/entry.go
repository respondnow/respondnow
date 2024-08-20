package server

import (
	"context"
	"github.com/respondnow/respond/server/pkg/auth"
	auth2 "github.com/respondnow/respond/server/pkg/database/mongodb/auth"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/respondnow/respond/server/api/middleware"
	"github.com/respondnow/respond/server/api/routes"
	"github.com/respondnow/respond/server/config"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	"github.com/respondnow/respond/server/pkg/prometheus"
	"github.com/respondnow/respond/server/utils"
	"github.com/sirupsen/logrus"
)

func Start() {
	if err := run(); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}

func run() error {
	loadConfig()

	logrus.Infof("go version %s", runtime.Version())
	logrus.Infof("go os %s", runtime.GOOS)
	logrus.Infof("go arch %s", runtime.GOARCH)

	// initialize mongo database
	err := mongodb.InitMongoClient()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, metricApp := setupServers()
	srv := startServer(ctx, app, config.EnvConfig.Ports.HttpPort, "HTTP server")
	metricSrv := startServer(ctx, metricApp, config.EnvConfig.Ports.MetricPort, "Metric server")

	backgroundProcess()

	<-ctx.Done()
	logrus.Infof("Received shutdown signal, shutting down servers...")

	shutdownServer(srv, "HTTP server")
	shutdownServer(metricSrv, "Metric server")

	return ctx.Err()
}

func loadConfig() {
	if err := envconfig.Process("", &config.EnvConfig); err != nil {
		logrus.Fatal(err)
	}

	if config.EnvConfig.Flags.EnableSLIMetrics {
		prometheus.Init()
	}
}

func setupServers() (*gin.Engine, *gin.Engine) {
	app := configureGin()
	metricApp := configureGin()
	registerRoutes(app, metricApp)
	return app, metricApp
}

func configureGin() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(middleware.DefaultStructuredLogger(), gin.Recovery())
	return app
}

func registerRoutes(app, metricApp *gin.Engine) {
	if config.EnvConfig.Flags.EnableSLIMetrics {
		app.Use(middleware.RequestMetricsMiddleware(), middleware.SLIAPIResponseTimeMiddleware())
	}

	routes.BaseRouter(app.Group("/"))
	routes.MetricRouter(metricApp.Group("/promMetrics"))
	routes.IncidentRouter(app.Group("/incident"))
	routes.AuthRouter(app.Group("/auth"))

	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, utils.DefaultResponseDTO{
			Status:  string(utils.ERROR),
			Message: "Page not found",
		})
	})

	logRegisteredRoutes(app)
	logRegisteredRoutes(metricApp)
}

func logRegisteredRoutes(engine *gin.Engine) {
	for _, route := range engine.Routes() {
		logrus.Infof("Registered route %s %s", route.Method, route.Path)
	}
}

func startServer(ctx context.Context, engine *gin.Engine, port, serverName string) *http.Server {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	go func() {
		logrus.Infof("%s starting on port: %s", serverName, port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start %s due to: %v", serverName, err)
		}
	}()

	return srv
}

func shutdownServer(srv *http.Server, serverName string) {
	logrus.Infof("Shutting down %s...", serverName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Failed to shutdown %s gracefully: %v", serverName, err)
	}
	logrus.Infof("%s shutdown complete", serverName)
}

func backgroundProcess() {
	go func() {
		_, err := auth.NewAuthService(auth2.NewAuthOperator(mongodb.Operator)).Signup(context.Background(), auth.AddUserInput{
			Name:     config.EnvConfig.Auth.DefaultUserName,
			UserID:   config.EnvConfig.Auth.DefaultUserID,
			Email:    config.EnvConfig.Auth.DefaultUserEmail,
			Password: config.EnvConfig.Auth.DefaultUserPassword,
		})
		if err != nil {
			logrus.Warnf("failed to create default user, err: %v", err)
		}

		logrus.Infof("Default user has been created successfully")
	}()
}

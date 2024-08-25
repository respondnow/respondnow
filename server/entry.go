package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/respondnow/respond/server/pkg/database/mongodb/hierarchy"
	hierarchy2 "github.com/respondnow/respond/server/pkg/hierarchy"
	"github.com/respondnow/respond/server/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"

	auth2 "github.com/respondnow/respond/server/pkg/database/mongodb/user"

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

const retryLimit = 3

func backgroundProcess() {
	go func() {
		ctx := context.Background()
		userService := user.NewAuthService(auth2.NewAuthOperator(mongodb.Operator))
		hierarchyService := hierarchy2.NewHierarchyManager(hierarchy.NewHierarchyOperator(mongodb.Operator))

		createdUser, err := createDefaultUser(ctx, userService)
		if err != nil {
			return
		}

		defaultAccount, err := createDefaultAccount(ctx, hierarchyService)
		if err != nil {
			fallbackCleanup(ctx, userService, hierarchyService, createdUser.ID, nil, nil, nil)
			return
		}

		defaultOrg, err := createDefaultOrganization(ctx, hierarchyService, defaultAccount)
		if err != nil {
			fallbackCleanup(ctx, userService, hierarchyService, createdUser.ID, &defaultAccount, nil, nil)
			return
		}

		defaultProject, err := createDefaultProject(ctx, hierarchyService, defaultAccount, defaultOrg)
		if err != nil {
			fallbackCleanup(ctx, userService, hierarchyService, createdUser.ID, &defaultAccount, &defaultOrg, nil)
			return
		}

		_, err = createUserMapping(ctx, hierarchyService, createdUser, defaultAccount, defaultOrg, defaultProject)
		if err != nil {
			fallbackCleanup(ctx, userService, hierarchyService, createdUser.ID, &defaultAccount, &defaultOrg, &defaultProject)
			return
		}

		logrus.Infof("All resources have been successfully created: User, Account, Organization, Project, and User Mapping.")
	}()
}

func createDefaultUser(ctx context.Context, userService user.AuthService) (auth2.User, error) {
	var createdUser auth2.User
	err := retry(retryLimit, func() error {
		var err error
		createdUser, err = userService.Signup(ctx, user.AddUserInput{
			Name:     config.EnvConfig.Auth.DefaultUserName,
			UserID:   config.EnvConfig.Auth.DefaultUserID,
			Email:    config.EnvConfig.Auth.DefaultUserEmail,
			Password: config.EnvConfig.Auth.DefaultUserPassword,
		})
		return err
	})
	if err != nil {
		logrus.Warnf("Failed to create default user after retries, err: %v", err)
	}
	return createdUser, err
}

func createDefaultAccount(ctx context.Context, hierarchyService hierarchy2.HierarchyManager) (hierarchy.Account, error) {
	defaultAccount := hierarchy.Account{
		ID:        primitive.NewObjectID(),
		AccountID: config.EnvConfig.DefaultHierarchy.DefaultAccountId,
		Name:      config.EnvConfig.DefaultHierarchy.DefaultAccountName,
		CreatedBy: config.SYSTEM,
		UpdatedBy: config.SYSTEM,
		CreatedAt: time.Now().Unix(),
	}
	err := retry(retryLimit, func() error {
		return hierarchyService.CreateAccount(ctx, defaultAccount)
	})
	if err != nil {
		logrus.Warnf("Failed to create default account after retries, err: %v", err)
	}
	return defaultAccount, err
}

func createDefaultOrganization(ctx context.Context, hierarchyService hierarchy2.HierarchyManager, defaultAccount hierarchy.Account) (hierarchy.Organization, error) {
	defaultOrg := hierarchy.Organization{
		ID:        primitive.NewObjectID(),
		OrgID:     config.EnvConfig.DefaultHierarchy.DefaultOrgId,
		Name:      config.EnvConfig.DefaultHierarchy.DefaultOrgName,
		AccountID: defaultAccount.AccountID,
		CreatedBy: config.SYSTEM,
		UpdatedBy: config.SYSTEM,
		CreatedAt: time.Now().Unix(),
	}
	err := retry(retryLimit, func() error {
		return hierarchyService.CreateOrganization(ctx, defaultOrg)
	})
	if err != nil {
		logrus.Warnf("Failed to create default organization after retries, err: %v", err)
	}
	return defaultOrg, err
}

func createDefaultProject(ctx context.Context, hierarchyService hierarchy2.HierarchyManager, defaultAccount hierarchy.Account, defaultOrg hierarchy.Organization) (hierarchy.Project, error) {
	defaultProject := hierarchy.Project{
		ID:        primitive.NewObjectID(),
		ProjectID: config.EnvConfig.DefaultHierarchy.DefaultProjectId,
		Name:      config.EnvConfig.DefaultHierarchy.DefaultProjectName,
		OrgID:     defaultOrg.OrgID,
		AccountID: defaultAccount.AccountID,
		CreatedBy: config.SYSTEM,
		UpdatedBy: config.SYSTEM,
		CreatedAt: time.Now().Unix(),
	}
	err := retry(retryLimit, func() error {
		return hierarchyService.CreateProject(ctx, defaultProject)
	})
	if err != nil {
		logrus.Warnf("Failed to create default project after retries, err: %v", err)
	}
	return defaultProject, err
}

func createUserMapping(ctx context.Context, hierarchyService hierarchy2.HierarchyManager, createdUser auth2.User, defaultAccount hierarchy.Account, defaultOrg hierarchy.Organization, defaultProject hierarchy.Project) (primitive.ObjectID, error) {
	var userMappingID primitive.ObjectID
	err := retry(retryLimit, func() error {
		var err error
		userMappingID, err = hierarchyService.CreateUserMapping(ctx, createdUser.UserID, defaultAccount.AccountID, defaultOrg.OrgID, defaultProject.ProjectID, true)
		return err
	})
	if err != nil {
		logrus.Warnf("Failed to create user mapping after retries, err: %v", err)
	}
	return userMappingID, err
}

func retry(limit int, fn func() error) error {
	var err error
	for i := 0; i < limit; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return err
}

func fallbackCleanup(ctx context.Context, userService user.AuthService, hierarchyService hierarchy2.HierarchyManager, id primitive.ObjectID, account *hierarchy.Account, org *hierarchy.Organization, project *hierarchy.Project) {
	logrus.Warnf("Cleaning up resources due to failure")

	if project != nil {
		err := hierarchyService.DeleteProject(ctx, project.ProjectID)
		if err != nil {
			logrus.Warnf("Failed to clean up project with ID: %v, err: %v", project.ProjectID, err)
		} else {
			logrus.Infof("Successfully cleaned up project with ID: %v", project.ProjectID)
		}
	}

	if org != nil {
		err := hierarchyService.DeleteOrganization(ctx, org.OrgID)
		if err != nil {
			logrus.Warnf("Failed to clean up organization with ID: %v, err: %v", org.OrgID, err)
		} else {
			logrus.Infof("Successfully cleaned up organization with ID: %v", org.OrgID)
		}
	}

	if account != nil {
		err := hierarchyService.DeleteAccount(ctx, account.AccountID)
		if err != nil {
			logrus.Warnf("Failed to clean up account with ID: %v, err: %v", account.AccountID, err)
		} else {
			logrus.Infof("Successfully cleaned up account with ID: %v", account.AccountID)
		}
	}

	err := userService.DeleteUser(ctx, id)
	if err != nil {
		logrus.Warnf("Failed to clean up user with ID: %v, err: %v", id, err)
	} else {
		logrus.Infof("Successfully cleaned up user with ID: %v", id)
	}
}

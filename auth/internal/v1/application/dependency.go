package application

import (
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/controller"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/repository"
	"github.com/PickHD/singkatin-revamp/auth/internal/v1/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
	AuthController        controller.AuthController
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.DB)
	authRepoImpl := repository.NewAuthRepository(app.Context, app.Config, app.Logger, app.DB)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, healthCheckRepoImpl)
	authSvcImpl := service.NewAuthService(app.Context, app.Config, app.Logger, authRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, healthCheckSvcImpl)
	authControllerImpl := controller.NewAuthController(app.Context, app.Config, app.Logger, authSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
		AuthController:        authControllerImpl,
	}
}

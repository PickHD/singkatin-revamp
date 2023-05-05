package application

import (
	"github.com/PickHD/singkatin-revamp/user/internal/v1/controller"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/repository"
	"github.com/PickHD/singkatin-revamp/user/internal/v1/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
	UserController        controller.UserController
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.DB)
	userRepoImpl := repository.NewUserRepository(app.Context, app.Config, app.Logger, app.DB)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, healthCheckRepoImpl)
	userSvcImpl := service.NewUserService(app.Context, app.Config, app.Logger, userRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, healthCheckSvcImpl)
	userControllerImpl := controller.NewUserController(app.Context, app.Config, app.Logger, userSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
		UserController:        userControllerImpl,
	}
}

package application

import (
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/controller"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/repository"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.DB, app.Redis)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, healthCheckRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, healthCheckSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
	}
}

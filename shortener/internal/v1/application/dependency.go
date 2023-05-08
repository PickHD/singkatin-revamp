package application

import (
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/controller"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/repository"
	"github.com/PickHD/singkatin-revamp/shortener/internal/v1/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
	ShortController       *controller.ShortControllerImpl
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.DB, app.Redis)
	shortRepoImpl := repository.NewShortRepository(app.Context, app.Config, app.Logger, app.DB, app.Redis, app.RabbitMQ)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, healthCheckRepoImpl)
	shortSvcImpl := service.NewShortService(app.Context, app.Config, app.Logger, shortRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, healthCheckSvcImpl)
	shortControllerImpl := controller.NewShortController(app.Context, app.Config, app.Logger, shortSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
		ShortController:       shortControllerImpl,
	}
}

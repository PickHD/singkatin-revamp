package application

import (
	"github.com/PickHD/singkatin-revamp/upload/internal/v1/controller"
	"github.com/PickHD/singkatin-revamp/upload/internal/v1/repository"
	"github.com/PickHD/singkatin-revamp/upload/internal/v1/service"
)

type Dependency struct {
	UploadController controller.UploadController
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	uploadRepo := repository.NewUploadRepository(app.Context, app.Config, app.Logger, app.Tracer, app.MinIO)

	// service
	uploadSvc := service.NewUploadService(app.Context, app.Config, app.Logger, app.Tracer, uploadRepo)

	// controller
	uploadController := controller.NewUploadController(app.Context, app.Config, app.Logger, app.Tracer, uploadSvc)

	return &Dependency{
		UploadController: uploadController,
	}
}

package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/internal/handlers"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/staff"
)

// SetupRoutes is the Router for GoFiber App
func (s *Server) SetupRoutes(app *fiber.App) {

	// Prepare a static middleware to serve the built React files.
	app.Static("/", "./web/build")

	// API routes group
	groupApiV1 := app.Group("/api/v:version?", handlers.ApiLimiter)
	{
		groupApiV1.Get("/", handlers.Index())
	}

	// App Repository
	staffRepo := staff.NewStaffRepository(s.MainDbConn)

	// auto migrate DB only on main process
	if !fiber.IsChild() {
		s.MainDbConn.AutoMigrate(models.Staff{})
	}
	// App Services
	staffService := staff.NewStaffService(staffRepo)
	// App Routes
	staff.NewStaffHandler(groupApiV1.Group("/staff"), staffService)

	// Prepare a fallback route to always serve the 'index.html', had there not be any matching routes.
	app.Static("*", "./web/build/index.html")
}

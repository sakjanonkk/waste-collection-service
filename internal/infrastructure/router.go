package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/internal/handlers"
	"github.com/zercle/gofiber-skelton/pkg/auth"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/staff"
	"github.com/zercle/gofiber-skelton/pkg/vehicle"

	"github.com/gofiber/swagger"
	_ "github.com/zercle/gofiber-skelton/docs"
)

// SetupRoutes is the Router for GoFiber App
func (s *Server) SetupRoutes(app *fiber.App) {

	app.Static("/", "./web/build")

	// API routes group
	groupApiV1 := app.Group("/api/v:version?", handlers.ApiLimiter)
	{
		groupApiV1.Get("/", handlers.Index())
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	// auto migrate DB only on main process
	if !fiber.IsChild() {
		s.MainDbConn.AutoMigrate(
			&models.Staff{},
			&models.Vehicle{},
		)
	}

	// Repositories
	staffRepo := staff.NewStaffRepository(s.MainDbConn)
	vehicleRepo := vehicle.NewVehicleRepository(s.MainDbConn)

	// Services
	staffService := staff.NewStaffService(staffRepo)
	vehicleService := vehicle.NewVehicleService(vehicleRepo)
	authService := auth.NewAuthService(staffRepo, s.JwtResources)

	// 🔓 Auth Routes (Public - No Authentication)
	authGroup := groupApiV1.Group("/auth")
	auth.NewAuthHandler(authGroup, authService, s.JwtResources)

	// 🔐 Staff Routes (Protected - Authentication Required)
	staffGroup := groupApiV1.Group("/staff")
	staffGroup.Use(auth.AuthMiddleware(s.JwtResources))
	staff.NewStaffHandler(staffGroup, staffService)

	// 🔐 Vehicle Routes (Protected - Authentication Required)
	vehicleGroup := groupApiV1.Group("/vehicles")
	vehicleGroup.Use(auth.AuthMiddleware(s.JwtResources)) // Require login for all vehicle routes
	vehicle.NewVehicleHandler(vehicleGroup, vehicleService)

	app.Static("*", "./web/build/index.html")
}

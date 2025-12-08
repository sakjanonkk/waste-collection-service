package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/internal/handlers"
	"github.com/zercle/gofiber-skelton/pkg/auth"
	"github.com/zercle/gofiber-skelton/pkg/collection_point"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/route"
	"github.com/zercle/gofiber-skelton/pkg/staff"
	"github.com/zercle/gofiber-skelton/pkg/vehicle"

	"github.com/gofiber/swagger"
	_ "github.com/zercle/gofiber-skelton/docs"
)

// SetupRoutes is the Router for GoFiber App
func (s *Server) SetupRoutes(app *fiber.App) {

	app.Static("/", "./web/build")
	app.Static("/uploads", "./uploads")

	// API routes group
	groupApiV1 := app.Group("/api/v:version?", handlers.ApiLimiter)
	{
		groupApiV1.Get("/", handlers.Index())
	}

	app.Get("/api/v1/swagger/*", swagger.HandlerDefault)

	// auto migrate DB only on main process
	if !fiber.IsChild() {
		s.MainDbConn.AutoMigrate(
			&models.Staff{},
			&models.Vehicle{},
			&models.CollectionPoint{},
			&models.Route{},
			&models.RoutePoint{},
		)
		SeedDefaultAdmin(s.MainDbConn)
		SeedTestData(s.MainDbConn)
	}

	// Repositories
	staffRepo := staff.NewStaffRepository(s.MainDbConn)
	vehicleRepo := vehicle.NewVehicleRepository(s.MainDbConn)
	collectionPointRepo := collection_point.NewCollectionPointRepository(s.MainDbConn)
	routeRepo := route.NewRouteRepository(s.MainDbConn)

	// Services
	staffService := staff.NewStaffService(staffRepo)
	vehicleService := vehicle.NewVehicleService(vehicleRepo)
	collectionPointService := collection_point.NewCollectionPointService(collectionPointRepo)
	routeService := route.NewRouteService(routeRepo)
	authService := auth.NewAuthService(staffRepo, s.JwtResources)
	// groupApiV1.Get("/hello-world", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })
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

	// 🔐 Collection Point Routes (Protected - Authentication Required)
	collectionPointGroup := groupApiV1.Group("/collection-points")
	collectionPointGroup.Use(auth.AuthMiddleware(s.JwtResources))
	collection_point.NewCollectionPointHandler(collectionPointGroup, collectionPointService)

	// 🔐 Route Routes (Protected - Authentication Required)
	routeGroup := groupApiV1.Group("/routes")
	routeGroup.Use(auth.AuthMiddleware(s.JwtResources))
	route.NewRouteHandler(routeGroup, routeService)

	app.Static("*", "./web/build/index.html")
}

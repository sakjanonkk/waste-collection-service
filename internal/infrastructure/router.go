package infrastructure

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/internal/handlers"
	"github.com/zercle/gofiber-skelton/pkg/auth"
	"github.com/zercle/gofiber-skelton/pkg/collection_point"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/permission"
	"github.com/zercle/gofiber-skelton/pkg/request"
	"github.com/zercle/gofiber-skelton/pkg/role"
	"github.com/zercle/gofiber-skelton/pkg/role_permission"
	"github.com/zercle/gofiber-skelton/pkg/route"
	"github.com/zercle/gofiber-skelton/pkg/staff"
	"github.com/zercle/gofiber-skelton/pkg/user_role"
	"github.com/zercle/gofiber-skelton/pkg/utils"
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
			&models.Request{},
		)
		SeedDefaultAdmin(s.MainDbConn)
		SeedTestData(s.MainDbConn)
		if err := models.MigrateTablePermissions(s.MainDbConn); err != nil {
			log.Panicf("Failed to migrate permissions: %v", err)
		}
		if err := utils.MigratePermission(s.MainDbConn); err != nil {
			log.Panicf("Failed to migrate permissions: %v", err)
		}
	}
	routerResources := handlers.NewRouterResources(s.JwtResources.JwtKeyfunc, s.MainDbConn)

	// Repositories
	staffRepo := staff.NewStaffRepository(s.MainDbConn)
	vehicleRepo := vehicle.NewVehicleRepository(s.MainDbConn)
	collectionPointRepo := collection_point.NewCollectionPointRepository(s.MainDbConn)
	routeRepo := route.NewRouteRepository(s.MainDbConn)
	requestRepo := request.NewRequestRepository(s.MainDbConn)
	permissionRepository := permission.NewPermissionRepository(s.Resources)
	roleRepository := role.NewRoleRepository(s.Resources)
	rolePermissionRepository := role_permission.NewRolePermissionRepository(s.Resources)
	userRoleRepository := user_role.NewUserRoleRepository(s.Resources)

	// Services
	staffService := staff.NewStaffService(staffRepo)
	vehicleService := vehicle.NewVehicleService(vehicleRepo)
	collectionPointService := collection_point.NewCollectionPointService(collectionPointRepo)
	routeService := route.NewRouteService(routeRepo)
	requestService := request.NewRequestService(requestRepo)
	authService := auth.NewAuthService(staffRepo, s.JwtResources)
	permissionService := permission.NewPermissionService(permissionRepository)
	roleService := role.NewRoleService(roleRepository)
	rolePermissionService := role_permission.NewRolePermissionService(rolePermissionRepository)
	userRoleService := user_role.NewUserRoleService(userRoleRepository)
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

	// 🔓 Request Routes (Public - Create Request)
	requestGroupPublic := groupApiV1.Group("/requests")
	request.NewRequestPublicHandler(requestGroupPublic, requestService)

	// 🔐 Request Routes (Protected - Approve/Reject/Get)
	requestGroup := groupApiV1.Group("/requests")
	requestGroup.Use(auth.AuthMiddleware(s.JwtResources))
	request.NewRequestProtectedHandler(requestGroup, requestService)
	permission.NewPermissionHandler(groupApiV1.Group("/permissions"), routerResources, permissionService)
	role.NewRoleHandler(groupApiV1.Group("/roles"), routerResources, roleService)
	role_permission.NewRolePermissionHandler(groupApiV1.Group("/role-permissions"), routerResources, rolePermissionService)
	user_role.NewUserRoleHandler(groupApiV1.Group("/user-roles"), routerResources, userRoleService)

	app.Static("*", "./web/build/index.html")
}

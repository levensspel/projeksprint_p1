package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	authHandler "github.com/levensspel/go-gin-template/handler/auth"
	employeeHandler "github.com/levensspel/go-gin-template/handler/employee"
	fileHandler "github.com/levensspel/go-gin-template/handler/file"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	employeeRepository "github.com/levensspel/go-gin-template/repository/employee"
	fileRepository "github.com/levensspel/go-gin-template/repository/file"
	userRepository "github.com/levensspel/go-gin-template/repository/user"
	employeeService "github.com/levensspel/go-gin-template/service/employee"
	fileService "github.com/levensspel/go-gin-template/service/file"
	userService "github.com/levensspel/go-gin-template/service/user"

	_ "github.com/levensspel/go-gin-template/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(r *gin.Engine, db *pgxpool.Pool) {
	logger := logger.NewlogHandler()

	// api := r.Group("/v1")
	// {
	// 	// untuk memanfaatkan api versioning, uncomment dan pakai ini
	// }

	userRepo := userRepository.NewUserRepository(db)
	fileRepo := fileRepository.NewFileRepository(db)

	userService := userService.NewUserService(userRepo, logger)
	fileService := fileService.NewFileService(fileRepo, logger)

	userHdlr := userHandler.NewUserHandler(userService, logger)
	authHandler := authHandler.NewHandler(userService, logger)
	fileHandler := fileHandler.NewHandler(fileService, logger)

	employeeRepo := employeeRepository.NewEmployeeRepository(db)
	employeeService := employeeService.NewEmployeeService(employeeRepo, logger)
	employeeHdlr := employeeHandler.NewEmployeeHandler(employeeService, logger)

	swaggerRoute := r.Group("/")
	{
		//Route untuk Swagger
		swaggerRoute.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	controllers := r.Group("/v1")
	{
		auth := controllers.Group("/auth")
		{
			auth.POST("", authHandler.Post)
		}

		file := controllers.Group("/file")
		{
			file.POST("", fileHandler.Upload)
		}

		user := controllers.Group("/user")
		{
			user.GET("", middleware.Authorization, userHdlr.GetProfile)
			user.PUT("", middleware.Authorization, userHdlr.Update)
			user.DELETE("", middleware.Authorization, userHdlr.Delete)
		}

		employee := controllers.Group("/employee")
		{
			employee.GET("", middleware.Authorization, employeeHdlr.GetEmployees)
			employee.GET("/with-join", middleware.Authorization, employeeHdlr.GetEmployeesWithJoin)
		}
		// tambah route lainnya disini
	}

}

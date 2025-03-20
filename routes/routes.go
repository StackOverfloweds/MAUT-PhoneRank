package routes

import (
	"log"
	"os"

	auth "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/Auth"
	brand "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/Brand"
	Smartphone "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/Smartphone"
	user "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/User"
	exports "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/export_smartphone"
	"github.com/StackOverfloweds/MAUT-PhoneRank/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App) {
	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ALLOW"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Ensure the static files directory exists
	if _, err := os.Stat("./data/Smartphone_Scrapping"); os.IsNotExist(err) {
		log.Fatal("Static files directory does not exist: ./data/Smartphone_Scrapping")
	}
	app.Static("/static", "./data/Smartphone_Scrapping")

	// routes for authentication
	authRoutes := app.Group("/auth")
	authRoutes.Post("/register", auth.Register)
	authRoutes.Post("/login", auth.Login)
	authRoutes.Post("/verify-otp", auth.VerifyOTP)
	authRoutes.Post("/logout", auth.Logout)

	//make routes for profile
	userProf := app.Group("/user", middleware.JWTMiddleware())
	userProf.Put("/profile", user.UpdateProfile)

	// routes for smartphone
	smartphoneRoutes := app.Group("/smartphone")
	smartphoneRoutes.Get("/", Smartphone.GetAllSmartphones)
	smartphoneRoutes.Get("/:id", Smartphone.GetSmartphoneDetail)
	smartphoneRoutes.Post("/search-maut", Smartphone.SearchSmartphoneWithMAUT)
	// Routes for brands
	brandRoutes := app.Group("/brands")
	brandRoutes.Get("/name", brand.GetAllBrand)

	exportRoutes := app.Group("/export")
	exportRoutes.Get("/smartphone", exports.ExportJSONSmartphone)

}

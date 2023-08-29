package main

import (
	"net/http"

	"github.com/caiostarke/go-expert-apo/configs"
	_ "github.com/caiostarke/go-expert-apo/docs"
	"github.com/caiostarke/go-expert-apo/internal/entity"
	"github.com/caiostarke/go-expert-apo/internal/infra/database"
	"github.com/caiostarke/go-expert-apo/internal/webserver/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title	Go Expert API Example
// @version	1.0
// @description	Product API with authentication
// @termsOfService	http://swagger.io/terms/

// @contact.name	Caio Starke

// @host	localhost:8000
// @BasePath	/
// @SecurityDefinitions.apikey	ApiKeyAuth
// @in	header
// @name	Authorization
func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	productHandler := handler.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handler.NewUserHandler(userDB, configs.TokenAuth, configs.JwtExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	// user handlers
	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}

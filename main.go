package main

import (
	"accounts"
	"cart"
	"fmt"
	"log"
	"net/http"
	"products"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	govalidator.SetFieldsRequiredByDefault(true)
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	m, err := migrate.New(
		"file://db/migrate",
		"postgres://prime:belindat2014@localhost:5432/avy",
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		fmt.Println(err)
	}

	r := chi.NewRouter()

	r.Use(accounts.Authenticator)
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		Debug:            true,
	}))

	// map the routes to the specific handlers
	r.Post("/v1/api/users/signup", accounts.SignupController)
	r.Post("/v1/api/users/login", accounts.LoginController)
	r.Post("/v1/api/products/write", products.ProductPostController)
	r.Post("/v1/api/cart/write", cart.CartPostController)

	fmt.Println("Starting server at port :8080")
	http.ListenAndServe(":8080", r)
}

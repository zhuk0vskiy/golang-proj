package main

import (
	"backend/src/config"
	"backend/src/internal/controller"
	"backend/src/internal/controller/v1"
	"backend/src/internal/app"
	"backend/src/pkg/logger"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"backend/src/pkg/mongodb"
)

func main() {
	ctx := context.Background()
	fmt.Println("trying to read config")
	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf(c.JwtKey)
	// Create logger
	fmt.Println("trying to create logger")
	loggerFile, err := os.OpenFile(
		c.Logger.File,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func(loggerFile *os.File) {
		err := loggerFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(loggerFile)

	l := logger.New(c.Logger.Level, loggerFile)

	tokenAuth := jwtauth.New("HS256", []byte(c.JwtKey), nil)
	fmt.Println(c.JwtKey)

	fmt.Printf("trying to connect db %s with user %s . port: %s\n", 
		c.Database.Postgres.Database, 
		c.Database.Postgres.User,
		fmt.Sprintf("%s:%d", c.Database.Postgres.Host, c.Database.Postgres.Port),
	)
	db, err := newConn(ctx, &c.Database)
	if err != nil {
		l.Fatalf("failed to connect to database: %v", err)
	}

	fmt.Printf("trying to connect db %s", c.Database.MongoDB.URI)
	mdb, err := mongodb.New(
		c.Database.MongoDB.URI,
		c.Database.MongoDB.Database,
		c.Database.MongoDB.Bucket,
	)
	if err != nil {
		l.Fatalf("failed to connect to database: %v", err)
	}

	fmt.Println("trying to new app")
	a := app.NewApp(db, mdb, c, l)

	fmt.Println("trying to make controller")
	mux := chi.NewMux()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	mux.Use(middleware.Logger)

	mux.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/login", v1.LoginHandler(a))
			r.Post("/signin", v1.SignInHandler(a))



			
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(jwtauth.Authenticator(tokenAuth))
				r.Use(controller.ValidateUserRoleJWT)

				r.Get("/validation", v1.ValidationHandler(a))
			})

			r.Route("/test", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Get("/studios/{id}", v1.GetStudioHandler(a))
					r.Post("/photo", v1.AddPhotoHandler(a))
					r.Get("/photo", v1.GetPhotoHandler(a))
					r.Delete("/photo", v1.DeletePhotoHandler(a))
				})
			})

			r.Route("/reserves", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Post("/", v1.AddReserveHandler(a))
				})
				r.Delete("/{id}", v1.DeleteReserveHandler(a))
			})

			r.Route("/studios", func(r chi.Router) {

				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateUserRoleJWT)

					r.Get("/{id}", v1.GetStudioHandler(a))
					r.Get("/{id}/rooms", v1.GetRoomsByStudioHandler(a))
					r.Get("/{id}/producers", v1.GetProducerHandler(a))
					r.Get("/{id}/instrumentalists", v1.GetInstrumentalistHandler(a))
					r.Get("/{id}/equipments", v1.GetEquipmentHandler(a))

				})
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateAdminRoleJWT)

					r.Patch("/{id}", v1.UpdateStudioHandler(a))
					r.Delete("/{id}", v1.DeleteStudioHandler(a))
					r.Post("/", v1.AddStudioHandler(a))
				})

			})

			r.Route("/rooms", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateUserRoleJWT)

					r.Get("/{id}", v1.GetRoomHandler(a))
				})

				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateAdminRoleJWT)

					r.Post("/", v1.AddRoomHandler(a))
					r.Patch("/{id}", v1.UpdateRoomHandler(a))
					r.Delete("/{id}", v1.DeleteRoomHandler(a))
				})
			})

			r.Route("/producers", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateUserRoleJWT)

					r.Get("/{id}", v1.GetProducerHandler(a))
				})

				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateAdminRoleJWT)

					r.Post("/", v1.AddProducerHandler(a))
					r.Patch("/{id}", v1.UpdateProducerHandler(a))
					r.Delete("/{id}", v1.DeleteProducerHandler(a))
				})
			})

			r.Route("/instrumentalists", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateUserRoleJWT)

					r.Get("/{id}", v1.GetInstrumentalistHandler(a))
				})
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateAdminRoleJWT)

					r.Post("/", v1.AddInstrumentalistHandler(a))
					r.Patch("/{id}", v1.UpdateInstrumentalistHandler(a))
					r.Delete("/{id}", v1.DeleteInstrumentalistHandler(a))
				})
			})

			r.Route("/equipments", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateUserRoleJWT)

					r.Get("/{id}", v1.GetEquipmentHandler(a))
				})
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateAdminRoleJWT)

					r.Post("/", v1.AddEquipmentHandler(a))
					r.Patch("/{id}", v1.UpdateEquipmentHandler(a))
					r.Delete("/{id}", v1.DeleteEquipmentHandler(a))
				})
			})

			r.Route("/user", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwtauth.Authenticator(tokenAuth))
					r.Use(controller.ValidateAdminRoleJWT)
					r.Get("/reserves", v1.GetUserReservesHandler(a))
				})
			})

		})
	})

	go func() {
		metricsAddress := fmt.Sprintf("%s:%d", c.Prometheus.MetricHost, c.Prometheus.MetricPort)

		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())

		fmt.Printf("сервер метрик прослушивает адрес: %s\n", metricsAddress)
		http.ListenAndServe(metricsAddress, metricsMux)
	}()

	serverPort := fmt.Sprintf(":%s", c.HTTP.Port)
	fmt.Printf("server has started at port %s\n", serverPort)
	err = http.ListenAndServe(serverPort, mux)
	if err != nil {
		log.Fatal(err)
	}

	//tui.Run(db, c, l)
}

func newConn(ctx context.Context, cfg *config.DatabaseConfig) (pool *pgxpool.Pool, err error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		cfg.Postgres.Driver,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("подключение к БД: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("пинг БД: %w", err)
	}

	return pool, nil
}

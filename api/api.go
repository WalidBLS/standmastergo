package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"standmaster/api/controller"
	"standmaster/internal/interaction"
	"standmaster/internal/kermesse"
	"standmaster/internal/stand"
	"standmaster/internal/ticket"
	"standmaster/internal/tombola"
	"standmaster/internal/user"
	"standmaster/third_party/resend"
)

type APIServer struct {
	address string
	db      *sqlx.DB
}

func NewAPIServer(address string, db *sqlx.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()

	resendService := resend.NewResendService(os.Getenv("RESEND_API_KEY"), os.Getenv("RESEND_FROM_EMAIL"))

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	userRepository := user.NewRepository(s.db)
	userService := user.NewService(userRepository, resendService)
	userController := controller.NewUserController(userService, userRepository)
	userController.RegisterRoutes(router)

	standRepository := stand.NewRepository(s.db)
	standService := stand.NewService(standRepository)
	standController := controller.NewStandController(standService, userRepository)
	standController.RegisterRoutes(router)

	kermesseRepository := kermesse.NewRepository(s.db)
	kermesseService := kermesse.NewService(kermesseRepository, userRepository)
	kermesseController := controller.NewKermesseController(kermesseService, userRepository)
	kermesseController.RegisterRoutes(router)

	interactionRepository := interaction.NewRepository(s.db)
	interactionService := interaction.NewService(interactionRepository, standRepository, userRepository, kermesseRepository)
	interactionController := controller.NewInteractionController(interactionService, userRepository)
	interactionController.RegisterRoutes(router)

	tombolaRepository := tombola.NewRepository(s.db)
	tombolaService := tombola.NewService(tombolaRepository, kermesseRepository)
	tombolaController := controller.NewTombolaController(tombolaService, userRepository)
	tombolaController.RegisterRoutes(router)

	ticketRepository := ticket.NewRepository(s.db)
	ticketService := ticket.NewService(ticketRepository, tombolaRepository, userRepository)
	ticketController := controller.NewTicketController(ticketService, userRepository)
	ticketController.RegisterRoutes(router)

	router.HandleFunc("/webhook", controller.HandleWebhook(userService)).Methods(http.MethodPost)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r := c.Handler(router)

	log.Printf("Starting server on %s", s.address)
	return http.ListenAndServe(s.address, r)
}

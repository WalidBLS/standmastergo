package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"standmaster/api/middleware"
	"standmaster/internal/models"
	"standmaster/internal/ticket"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/json"
)

type TicketController struct {
	service        ticket.TicketService
	userRepository user.UserRepository
}

func NewTicketController(service ticket.TicketService, userRepository user.UserRepository) *TicketController {
	return &TicketController{
		service:        service,
		userRepository: userRepository,
	}
}

func (h *TicketController) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/tickets", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userRepository, models.UserRoleOrganizer, models.UserRoleParent, models.UserRoleChild))).Methods(http.MethodGet)
	mux.Handle("/ticket/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userRepository, models.UserRoleOrganizer, models.UserRoleParent, models.UserRoleChild))).Methods(http.MethodGet)
	mux.Handle("/ticket", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userRepository, models.UserRoleChild))).Methods(http.MethodPost)
}

func (h *TicketController) GetAll(w http.ResponseWriter, r *http.Request) error {
	tickets, err := h.service.GetAll(r.Context())
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, tickets); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *TicketController) Get(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	ticket, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, ticket); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *TicketController) Create(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Create(r.Context(), input); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusCreated, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"standmaster/api/middleware"
	"standmaster/internal/models"
	"standmaster/internal/stand"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/json"
	"standmaster/pkg/utils"
)

type StandController struct {
	service        stand.StandService
	userRepository user.UserRepository
}

func NewStandController(service stand.StandService, userRepository user.UserRepository) *StandController {
	return &StandController{
		service:        service,
		userRepository: userRepository,
	}
}

func (h *StandController) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/stands", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userRepository))).Methods(http.MethodGet)
	mux.Handle("/stand/current", errors.ErrorHandler(middleware.IsAuth(h.GetCurrent, h.userRepository, models.UserRoleStandHolder))).Methods(http.MethodGet)
	mux.Handle("/stand/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userRepository))).Methods(http.MethodGet)
	mux.Handle("/stand", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userRepository, models.UserRoleStandHolder))).Methods(http.MethodPost)
	mux.Handle("/stand", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userRepository, models.UserRoleStandHolder))).Methods(http.MethodPatch)
}

func (h *StandController) GetAll(w http.ResponseWriter, r *http.Request) error {
	stands, err := h.service.GetAll(r.Context(), utils.GetQueryParams(r))
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, stands); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *StandController) Get(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	stand, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, stand); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *StandController) GetCurrent(w http.ResponseWriter, r *http.Request) error {
	stand, err := h.service.GetCurrent(r.Context())
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, stand); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *StandController) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *StandController) Update(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.UpdateCurrent(r.Context(), input); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusAccepted, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

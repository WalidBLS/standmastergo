package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"standmaster/api/middleware"
	"standmaster/internal/models"
	"standmaster/internal/tombola"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/json"
	"standmaster/pkg/utils"
)

type TombolaController struct {
	service        tombola.TombolaService
	userRepository user.UserRepository
}

func NewTombolaController(service tombola.TombolaService, userRepository user.UserRepository) *TombolaController {
	return &TombolaController{
		service:        service,
		userRepository: userRepository,
	}
}

func (h *TombolaController) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/tombolas", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userRepository))).Methods(http.MethodGet)
	mux.Handle("/tombola/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userRepository))).Methods(http.MethodGet)
	mux.Handle("/tombola", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userRepository, models.UserRoleOrganizer))).Methods(http.MethodPost)
	mux.Handle("/tombola/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userRepository, models.UserRoleOrganizer))).Methods(http.MethodPatch)
	mux.Handle("/tombola/{id}/finish", errors.ErrorHandler(middleware.IsAuth(h.Finish, h.userRepository, models.UserRoleOrganizer))).Methods(http.MethodPatch)
}

func (h *TombolaController) GetAll(w http.ResponseWriter, r *http.Request) error {
	tombolas, err := h.service.GetAll(r.Context(), utils.GetQueryParams(r))
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, tombolas); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *TombolaController) Get(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	tombola, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, tombola); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *TombolaController) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *TombolaController) Update(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Update(r.Context(), id, input); err != nil {
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

func (h *TombolaController) Finish(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Finish(r.Context(), id); err != nil {
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

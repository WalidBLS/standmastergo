package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"standmaster/api/middleware"
	"standmaster/internal/kermesse"
	"standmaster/internal/models"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/json"
)

type KermesseController struct {
	service        kermesse.KermesseService
	userRepository user.UserRepository
}

func NewKermesseController(service kermesse.KermesseService, userRepository user.UserRepository) *KermesseController {
	return &KermesseController{
		service:        service,
		userRepository: userRepository,
	}
}

func (h *KermesseController) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/kermesse", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userRepository, models.UserRoleOrganizer))).Methods(http.MethodPost)
	mux.Handle("/kermesse/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userRepository, models.UserRoleOrganizer))).Methods(http.MethodPatch)
	mux.Handle("/kermesse/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userRepository))).Methods(http.MethodGet)
	mux.Handle("/kermesses", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userRepository))).Methods(http.MethodGet)
	mux.Handle("/kermesse/{id}/users", errors.ErrorHandler(middleware.IsAuth(h.GetUsersInvite, h.userRepository))).Methods(http.MethodGet)
	mux.Handle("/kermesse/{id}/finish", errors.ErrorHandler(middleware.IsAuth(h.End, h.userRepository, models.UserRoleOrganizer))).Methods(http.MethodPatch)
	mux.Handle("/kermesse/{id}/adduser", errors.ErrorHandler(middleware.IsAuth(h.AddUser, h.userRepository, models.UserRoleOrganizer))).Methods(http.MethodPatch)
	mux.Handle("/kermesse/{id}/addstand", errors.ErrorHandler(middleware.IsAuth(h.AddStand, h.userRepository, models.UserRoleOrganizer))).Methods(http.MethodPatch)
}

func (h *KermesseController) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *KermesseController) Update(w http.ResponseWriter, r *http.Request) error {
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

func (h *KermesseController) Get(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	kermesse, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, kermesse); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *KermesseController) GetAll(w http.ResponseWriter, r *http.Request) error {
	kermesses, err := h.service.GetAll(r.Context())
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, kermesses); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *KermesseController) GetUsersInvite(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	users, err := h.service.GetUsersInvite(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, users); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *KermesseController) End(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.End(r.Context(), id); err != nil {
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

func (h *KermesseController) AddUser(w http.ResponseWriter, r *http.Request) error {
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
	input["kermesse_id"] = id

	if err := h.service.AddUser(r.Context(), input); err != nil {
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

func (h *KermesseController) AddStand(w http.ResponseWriter, r *http.Request) error {
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
	input["kermesse_id"] = id

	if err := h.service.AddStand(r.Context(), input); err != nil {
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

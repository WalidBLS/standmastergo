package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"standmaster/api/middleware"
	"standmaster/internal/models"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/json"
	"standmaster/pkg/utils"
)

type UserController struct {
	service    user.UserService
	repository user.UserRepository
}

func NewUserController(service user.UserService, repository user.UserRepository) *UserController {
	return &UserController{
		service:    service,
		repository: repository,
	}
}

func (h *UserController) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/users", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.repository))).Methods(http.MethodGet)
	mux.Handle("/user/children", errors.ErrorHandler(middleware.IsAuth(h.GetAllChildren, h.repository, models.UserRoleParent))).Methods(http.MethodGet)
	mux.Handle("/user/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.repository))).Methods(http.MethodGet)
	mux.Handle("/user/invite", errors.ErrorHandler(middleware.IsAuth(h.Invite, h.repository, models.UserRoleParent))).Methods(http.MethodPost)
	mux.Handle("/user/pay", errors.ErrorHandler(middleware.IsAuth(h.Pay, h.repository, models.UserRoleParent))).Methods(http.MethodPatch)
	mux.Handle("/user/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.repository))).Methods(http.MethodPatch)

	mux.Handle("/register", errors.ErrorHandler(h.SignUp)).Methods(http.MethodPost)
	mux.Handle("/login", errors.ErrorHandler(h.SignIn)).Methods(http.MethodPost)
	mux.Handle("/profile", errors.ErrorHandler(middleware.IsAuth(h.GetMe, h.repository))).Methods(http.MethodGet)
}

func (h *UserController) GetAll(w http.ResponseWriter, r *http.Request) error {
	users, err := h.service.GetAll(r.Context(), utils.GetQueryParams(r))
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

func (h *UserController) GetAllChildren(w http.ResponseWriter, r *http.Request) error {
	users, err := h.service.GetAllChildren(r.Context(), utils.GetQueryParams(r))
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

func (h *UserController) Get(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, user); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserController) Update(w http.ResponseWriter, r *http.Request) error {
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

func (h *UserController) Invite(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Invite(r.Context(), input); err != nil {
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

func (h *UserController) Pay(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Pay(r.Context(), input); err != nil {
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

func (h *UserController) SignUp(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.SignUp(r.Context(), input); err != nil {
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

func (h *UserController) SignIn(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	response, err := h.service.SignIn(r.Context(), input)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, response); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserController) GetMe(w http.ResponseWriter, r *http.Request) error {
	response, err := h.service.GetMe(r.Context())
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, response); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

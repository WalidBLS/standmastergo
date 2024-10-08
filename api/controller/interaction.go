package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"standmaster/api/middleware"
	"standmaster/internal/interaction"
	"standmaster/internal/models"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/json"
	"standmaster/pkg/utils"
)

type InteractionController struct {
	service        interaction.InteractionService
	userRepository user.UserRepository
}

func NewInteractionController(service interaction.InteractionService, userRepository user.UserRepository) *InteractionController {
	return &InteractionController{
		service:        service,
		userRepository: userRepository,
	}
}

func (h *InteractionController) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/interaction", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userRepository, models.UserRoleParent, models.UserRoleChild))).Methods(http.MethodPost)
	mux.Handle("/interaction/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userRepository))).Methods(http.MethodGet)
	mux.Handle("/interaction/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userRepository, models.UserRoleStandHolder))).Methods(http.MethodPatch)
	mux.Handle("/interactions", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userRepository))).Methods(http.MethodGet)
}

func (h *InteractionController) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *InteractionController) Update(w http.ResponseWriter, r *http.Request) error {
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

func (h *InteractionController) GetAll(w http.ResponseWriter, r *http.Request) error {
	interactions, err := h.service.GetAll(r.Context(), utils.GetQueryParams(r))
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, interactions); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *InteractionController) Get(w http.ResponseWriter, r *http.Request) error {
	queryParams := mux.Vars(r)
	id, err := strconv.Atoi(queryParams["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	interaction, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, interaction); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

package ticket

import (
	"context"
	"database/sql"
	goErrors "errors"

	"standmaster/internal/models"
	"standmaster/internal/tombola"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/utils"
)

type TicketService interface {
	GetAll(ctx context.Context) ([]models.Ticket, error)
	Get(ctx context.Context, id int) (models.Ticket, error)
	Create(ctx context.Context, input map[string]interface{}) error
}

type Service struct {
	repository        TicketRepository
	tombolaRepository tombola.TombolaRepository
	userRepository    user.UserRepository
}

func NewService(repository TicketRepository, tombolaRepository tombola.TombolaRepository, userRepository user.UserRepository) *Service {
	return &Service{
		repository:        repository,
		tombolaRepository: tombolaRepository,
		userRepository:    userRepository,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]models.Ticket, error) {
	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return nil, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}
	userRole, ok := ctx.Value(models.UserRoleKey).(string)
	if !ok {
		return nil, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user role not found in context"),
		}
	}

	filters := map[string]interface{}{}
	if userRole == models.UserRoleOrganizer {
		filters["organizer_id"] = userId
	} else if userRole == models.UserRoleParent {
		filters["parent_id"] = userId
	} else if userRole == models.UserRoleChild {
		filters["child_id"] = userId
	}

	tickets, err := s.repository.FindAll(filters)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return tickets, nil
}

func (s *Service) Get(ctx context.Context, id int) (models.Ticket, error) {
	ticket, err := s.repository.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return ticket, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return ticket, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return ticket, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	tombolaId, err := utils.GetIntFromMap(input, "tombola_id")
	if err != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: err,
		}
	}
	tombola, err := s.tombolaRepository.FindById(tombolaId)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if tombola.Status != models.TombolaStatusStarted {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("tombola is not started or already finished"),
		}
	}

	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}
	user, err := s.userRepository.FindById(userId)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	// check if user has enough credit
	if user.Credit < tombola.Price {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("not enough credit"),
		}
	}

	// check if user belongs to the kermesse
	canCreate, err := s.repository.CanCreate(map[string]interface{}{
		"kermesse_id": tombola.KermesseId,
		"user_id":     userId,
	})
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	if !canCreate {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("forbidden"),
		}
	}

	// decrease user's credit
	err = s.userRepository.UpdateCredit(userId, -tombola.Price)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	input["user_id"] = userId

	err = s.repository.Create(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

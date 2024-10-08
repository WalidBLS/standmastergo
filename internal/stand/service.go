package stand

import (
	"context"
	"database/sql"
	goErrors "errors"

	"standmaster/internal/models"
	"standmaster/pkg/errors"
)

type StandService interface {
	GetAll(ctx context.Context, params map[string]interface{}) ([]models.Stand, error)
	Get(ctx context.Context, id int) (models.Stand, error)
	GetCurrent(ctx context.Context) (models.Stand, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	UpdateCurrent(ctx context.Context, input map[string]interface{}) error
}

type Service struct {
	repository StandRepository
}

func NewService(repository StandRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAll(ctx context.Context, params map[string]interface{}) ([]models.Stand, error) {
	filters := map[string]interface{}{}
	if params["kermesse_id"] != nil {
		filters["kermesse_id"] = params["kermesse_id"]
	}
	if params["is_free"] != nil {
		filters["is_free"] = params["is_free"]
	}

	stands, err := s.repository.FindAll(params)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return stands, nil
}

func (s *Service) Get(ctx context.Context, id int) (models.Stand, error) {
	stand, err := s.repository.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return stand, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return stand, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return stand, nil
}

func (s *Service) GetCurrent(ctx context.Context) (models.Stand, error) {
	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return models.Stand{}, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}

	stand, err := s.repository.FindByUserId(userId)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return stand, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return stand, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return stand, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}
	input["user_id"] = userId

	err := s.repository.Create(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id int, input map[string]interface{}) error {
	stand, err := s.repository.FindById(id)
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

	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}
	if stand.UserId != userId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("user is not the holder of the stand"),
		}
	}

	err = s.repository.Update(id, input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) UpdateCurrent(ctx context.Context, input map[string]interface{}) error {
	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}

	err := s.repository.UpdateByUserId(userId, input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

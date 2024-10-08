package tombola

import (
	"context"
	"database/sql"
	goErrors "errors"

	"standmaster/internal/kermesse"
	"standmaster/internal/models"
	"standmaster/pkg/errors"
	"standmaster/pkg/utils"
)

type TombolaService interface {
	GetAll(ctx context.Context, params map[string]interface{}) ([]models.Tombola, error)
	Get(ctx context.Context, id int) (models.Tombola, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	Finish(ctx context.Context, id int) error
}

type Service struct {
	repository         TombolaRepository
	kermesseRepository kermesse.KermesseRepository
}

func NewService(repository TombolaRepository, kermesseRepository kermesse.KermesseRepository) *Service {
	return &Service{
		repository:         repository,
		kermesseRepository: kermesseRepository,
	}
}

func (s *Service) GetAll(ctx context.Context, params map[string]interface{}) ([]models.Tombola, error) {
	filters := map[string]interface{}{}
	if params["kermesse_id"] != nil {
		filters["kermesse_id"] = params["kermesse_id"]
	}

	tombolas, err := s.repository.FindAll(filters)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return tombolas, nil
}

func (s *Service) Get(ctx context.Context, id int) (models.Tombola, error) {
	tombola, err := s.repository.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return tombola, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return tombola, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return tombola, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	kermesseId, error := utils.GetIntFromMap(input, "kermesse_id")
	if error != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: error,
		}
	}
	kermesse, err := s.kermesseRepository.FindById(kermesseId)
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

	if kermesse.Status == models.KermesseStatusEnded {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("kermesse is ended"),
		}
	}

	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}
	if kermesse.UserId != userId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("forbidden"),
		}
	}

	err = s.repository.Create(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id int, input map[string]interface{}) error {
	tombola, err := s.repository.FindById(id)
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

	kermesse, err := s.kermesseRepository.FindById(tombola.KermesseId)
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

	if kermesse.Status == models.KermesseStatusEnded {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("kermesse is ended"),
		}
	}

	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}
	if kermesse.UserId != userId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("forbidden"),
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

func (s *Service) Finish(ctx context.Context, id int) error {
	tombola, err := s.repository.FindById(id)
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

	kermesse, err := s.kermesseRepository.FindById(tombola.KermesseId)
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

	if kermesse.Status == models.KermesseStatusEnded {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("kermesse is ended"),
		}
	}

	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}
	if kermesse.UserId != userId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("forbidden"),
		}
	}

	if tombola.Status != models.TombolaStatusStarted {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("tombola is not started"),
		}
	}

	err = s.repository.SetWinner(id)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

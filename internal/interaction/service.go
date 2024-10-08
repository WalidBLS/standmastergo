package interaction

import (
	"context"
	"database/sql"
	goErrors "errors"

	"standmaster/internal/kermesse"
	"standmaster/internal/models"
	"standmaster/internal/stand"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/utils"
)

type InteractionService interface {
	GetAll(ctx context.Context, params map[string]interface{}) ([]models.InteractionBasic, error)
	Get(ctx context.Context, id int) (models.Interaction, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
}

type Service struct {
	repository         InteractionRepository
	standRepository    stand.StandRepository
	userRepository     user.UserRepository
	kermesseRepository kermesse.KermesseRepository
}

func NewService(repository InteractionRepository, standRepository stand.StandRepository, userRepository user.UserRepository, kermesseRepository kermesse.KermesseRepository) *Service {
	return &Service{
		repository:         repository,
		standRepository:    standRepository,
		userRepository:     userRepository,
		kermesseRepository: kermesseRepository,
	}
}

func (s *Service) GetAll(ctx context.Context, params map[string]interface{}) ([]models.InteractionBasic, error) {

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
	if userRole == models.UserRoleParent {
		filters["parent_id"] = userId
	} else if userRole == models.UserRoleChild {
		filters["child_id"] = userId
	} else if userRole == models.UserRoleStandHolder {
		filters["stand_holder_id"] = userId
	}
	if params["kermesse_id"] != nil {
		filters["kermesse_id"] = params["kermesse_id"]
	}

	interactions, err := s.repository.FindAll(filters)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return interactions, nil
}

func (s *Service) Get(ctx context.Context, id int) (models.Interaction, error) {
	interaction, err := s.repository.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return interaction, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return interaction, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return interaction, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	standId, err := utils.GetIntFromMap(input, "stand_id")
	if err != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: err,
		}
	}
	stand, err := s.standRepository.FindById(standId)
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

	canCreate, err := s.repository.CanCreate(map[string]interface{}{
		"user_id":  userId,
		"stand_id": standId,
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

	// calculate total price
	quantity := 1
	totalPrice := stand.Price
	if stand.Type == models.InteractionTypeConsumption {
		quantity, err = utils.GetIntFromMap(input, "quantity")
		if err != nil {
			return errors.CustomError{
				Key: errors.BadRequest,
				Err: err,
			}
		}
		totalPrice = stand.Price * quantity
	}

	// check stand's stock and user credit
	if stand.Type == models.InteractionTypeConsumption {
		if stand.Stock < quantity {
			return errors.CustomError{
				Key: errors.BadRequest,
				Err: goErrors.New("not enough stock"),
			}
		}
	}

	// check user's credit
	if user.Credit < totalPrice {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("not enough credit"),
		}
	}

	// decrease stand's stock
	if stand.Type == models.InteractionTypeConsumption {
		err = s.standRepository.UpdateStock(standId, -quantity)
		if err != nil {
			return errors.CustomError{
				Key: errors.InternalServerError,
				Err: err,
			}
		}
	}

	// decrease user's credit
	err = s.userRepository.UpdateCredit(userId, -totalPrice)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	// increase stand holder's credit
	err = s.userRepository.UpdateCredit(stand.UserId, totalPrice)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	input["user_id"] = user.Id
	input["type"] = stand.Type
	input["credit"] = totalPrice

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
	interaction, err := s.repository.FindById(id)
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

	if interaction.Type != models.InteractionTypeActivity {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("interaction type is not activity"),
		}
	}

	kermesse, err := s.kermesseRepository.FindById(interaction.Kermesse.Id)
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

	stand, err := s.standRepository.FindById(interaction.Stand.Id)
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
			Err: goErrors.New("forbidden"),
		}
	}

	err = s.repository.Update(id, map[string]interface{}{
		"status": models.InteractionStatusEnded,
		"point":  input["point"],
	})
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

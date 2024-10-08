package kermesse

import (
	"context"
	"database/sql"
	goErrors "errors"

	"standmaster/internal/models"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/utils"
)

type KermesseService interface {
	GetAll(ctx context.Context) ([]models.Kermesse, error)
	GetUsersInvite(ctx context.Context, id int) ([]models.UserBasic, error)
	Get(ctx context.Context, id int) (models.KermesseWithStats, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	End(ctx context.Context, id int) error

	AddUser(ctx context.Context, input map[string]interface{}) error
	AddStand(ctx context.Context, input map[string]interface{}) error
}

type Service struct {
	repository     KermesseRepository
	userRepository user.UserRepository
}

func NewService(repository KermesseRepository, userRepository user.UserRepository) *Service {
	return &Service{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]models.Kermesse, error) {
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
	} else if userRole == models.UserRoleStandHolder {
		filters["stand_holder_id"] = userId
	}

	kermesses, err := s.repository.FindAll(filters)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return kermesses, nil
}

func (s *Service) GetUsersInvite(ctx context.Context, id int) ([]models.UserBasic, error) {
	users, err := s.repository.FindUsersInvite(id)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return users, nil
}

func (s *Service) Get(ctx context.Context, id int) (models.KermesseWithStats, error) {
	userId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return models.KermesseWithStats{}, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}
	userRole, ok := ctx.Value(models.UserRoleKey).(string)
	if !ok {
		return models.KermesseWithStats{}, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user role not found in context"),
		}
	}

	kermesse, err := s.repository.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return models.KermesseWithStats{}, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return models.KermesseWithStats{}, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	filters := map[string]interface{}{}
	if userRole == models.UserRoleOrganizer {
		filters["organizer_id"] = userId
	} else if userRole == models.UserRoleParent {
		filters["parent_id"] = userId
	} else if userRole == models.UserRoleChild {
		filters["child_id"] = userId
	} else if userRole == models.UserRoleStandHolder {
		filters["stand_holder_id"] = userId
	}

	stats, err := s.repository.Stats(id, filters)
	if err != nil {
		return models.KermesseWithStats{}, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	kermesseWithStats := models.KermesseWithStats{
		Id:                kermesse.Id,
		UserId:            kermesse.UserId,
		Name:              kermesse.Name,
		Description:       kermesse.Description,
		Status:            kermesse.Status,
		StandCount:        stats.StandCount,
		TombolaCount:      stats.TombolaCount,
		UserCount:         stats.UserCount,
		InteractionCount:  stats.InteractionCount,
		InteractionIncome: stats.InteractionIncome,
		TombolaIncome:     stats.TombolaIncome,
	}

	return kermesseWithStats, nil
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
	kermesse, err := s.repository.FindById(id)
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
			Err: goErrors.New("kermesse is already ended"),
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

func (s *Service) End(ctx context.Context, id int) error {
	kermesse, err := s.repository.FindById(id)
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
			Err: goErrors.New("kermesse is already ended"),
		}
	}

	canEnd, err := s.repository.CanEnd(id)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	if !canEnd {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("kermesse can't be ended"),
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

	err = s.repository.End(id)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) AddUser(ctx context.Context, input map[string]interface{}) error {
	kermesse, err := s.repository.FindById(input["kermesse_id"].(int))
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
			Err: goErrors.New("kermesse is already ended"),
		}
	}

	organizerId, ok := ctx.Value(models.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}
	if kermesse.UserId != organizerId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("forbidden"),
		}
	}

	childId, error := utils.GetIntFromMap(input, "user_id")
	if error != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: error,
		}
	}
	child, err := s.userRepository.FindById(childId)
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
	if child.Role != models.UserRoleChild {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("user is not a child"),
		}
	}

	// invite child
	err = s.repository.AddUser(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	// invite child parent if exists
	if child.ParentId != nil {
		input["user_id"] = child.ParentId
		err = s.repository.AddUser(input)
		if err != nil {
			return errors.CustomError{
				Key: errors.InternalServerError,
				Err: err,
			}
		}
	}

	return nil
}

func (s *Service) AddStand(ctx context.Context, input map[string]interface{}) error {
	kermesse, err := s.repository.FindById(input["kermesse_id"].(int))
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
			Err: goErrors.New("kermesse is already ended"),
		}
	}

	standId, err := utils.GetIntFromMap(input, "stand_id")
	if err != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: err,
		}
	}
	canAddStand, err := s.repository.CanAddStand(standId)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	if !canAddStand {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("stand is already associated with kermesse"),
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

	err = s.repository.AddStand(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

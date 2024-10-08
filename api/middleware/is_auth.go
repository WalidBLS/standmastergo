package middleware

import (
	"context"
	goErrors "errors"
	"net/http"
	"os"
	"strings"

	"standmaster/internal/models"
	"standmaster/internal/user"
	"standmaster/pkg/errors"
	"standmaster/pkg/jwt"
)

func IsAuth(handlerFunc errors.ErrorHandler, repository user.UserRepository, roles ...string) errors.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		token := r.Header.Get("Authorization")
		if token == "" {
			return errors.CustomError{
				Key: errors.Unauthorized,
				Err: goErrors.New("token is required"),
			}
		}

		tokenParts := strings.Split(token, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return errors.CustomError{
				Key: errors.Unauthorized,
				Err: goErrors.New("invalid token"),
			}
		}

		userId, err := jwt.GetTokenUserId(tokenParts[1], os.Getenv("JWT_SECRET"))
		if err != nil {
			return errors.CustomError{
				Key: errors.Unauthorized,
				Err: goErrors.New("invalid token"),
			}
		}

		user, err := repository.FindById(userId)
		if err != nil {
			return err
		}

		if len(roles) > 0 {
			roleAllowed := false
			for _, role := range roles {
				if user.Role == role {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				return errors.CustomError{
					Key: errors.Forbidden,
					Err: goErrors.New("user does not have the required role"),
				}
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, models.UserIDKey, user.Id)
		ctx = context.WithValue(ctx, models.UserRoleKey, user.Role)
		r = r.WithContext(ctx)

		return handlerFunc(w, r)
	}
}

package middleware

import (
	"errors"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"guardian/internal/models/entities"
	"guardian/utlis/logger"
	"net/http"
	"strings"
	"time"

	"guardian/configs"

	"github.com/golang-jwt/jwt/v4"
)

func VerifyJWT(protected http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if configs.GlobalConfig.EnableExternalAuth {
			err := VerifyExternalJWT(r)
			if err == nil {
				protected.ServeHTTP(w, r)
				return
			}
		}

		jwtVerifier := jwtauth.Verifier(configs.GlobalConfig.TokenAuth)
		jwtAuthenticator := jwtauth.Authenticator(configs.GlobalConfig.TokenAuth)

		jwtVerifier(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtAuthenticator(protected).ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	})
}

func VerifyExternalJWT(r *http.Request) error {
	tokenStr, err := extractToken(r)
	if err != nil {
		return fmt.Errorf("JWT token not provided: %v", err)
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return configs.GlobalConfig.Jwk.Keyfunc(token)
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !validateClaims(claims) {
		return fmt.Errorf("invalid claims")
	}

	userIDStr := claims["user_id"].(string)
	_, err = checkUserInDB(userIDStr)
	if err != nil {
		_, err = createUserInDB(userIDStr)
		if err != nil {
			return fmt.Errorf("unable to register the user: %v", err)
		}
	}

	return nil
}

func validateClaims(claims jwt.MapClaims) bool {
	if claims["iss"] != configs.GlobalConfig.ExternalJwtIssuer {
		return false
	}

	if claims["aud"] != configs.GlobalConfig.ExternalJwtAudience {
		return false
	}

	if claims["exp"] != nil && claims["exp"].(float64) < float64(time.Now().Unix()) {
		return false
	}

	return true
}

func extractToken(r *http.Request) (string, error) {
	adminAuthHeader := r.Header.Get("X-Guardian-Authorization")
	if adminAuthHeader != "" {
		return adminAuthHeader, nil
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	return parts[1], nil
}

func checkUserInDB(userID string) (*entities.User, error) {
	return nil, nil
}

func createUserInDB(userID string) (*entities.User, error) {
	return nil, nil
}

func ParseRequestJWT(r *http.Request) (jwt.MapClaims, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve user claims: %v", err)
	}
	return claims, nil
}

func GetUserFromContext(r *http.Request) (*primitive.ObjectID, error) {
	claims, err := ParseRequestJWT(r)
	if err != nil {
		return nil, err
	}

	userIDStr := claims["user_id"].(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		logger.GetLogger().Errorf("error in reading the userID from the user's token: %s", userIDStr)
		return nil, err
	}

	return &userID, nil
}

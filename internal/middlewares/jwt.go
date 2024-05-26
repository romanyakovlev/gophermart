package middlewares

import (
	"context"
	"net/http"

	golangjwt "github.com/golang-jwt/jwt/v4"

	"github.com/romanyakovlev/gophermart/internal/jwt"
	"github.com/romanyakovlev/gophermart/internal/models"
)

type contextKey string

const userContextKey contextKey = "currentUser"

func GetUserFromContext(ctx context.Context) (models.AuthUser, bool) {
	user, ok := ctx.Value(userContextKey).(models.AuthUser)
	return user, ok
}

/*
TODO: создать MockedJWTMiddleware
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var authUser models.AuthUser
		uuidStr := "e0a308ae-4132-4fa5-a70c-5cf559ec8933"
		UUID, _ := uuid.Parse(uuidStr)
		authUser = models.AuthUser{UUID: UUID}
		ctxWithUser := context.WithValue(r.Context(), userContextKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}
*/

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/register" || r.URL.Path == "/api/user/login" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value
		claims := &models.Claims{}
		token, err := golangjwt.ParseWithClaims(tokenStr, claims, func(token *golangjwt.Token) (interface{}, error) {
			return jwt.JWTKey, nil
		})

		if err != nil {
			if err == golangjwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authUser := models.AuthUser{UUID: claims.UUID}

		ctxWithUser := context.WithValue(r.Context(), userContextKey, authUser)

		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}

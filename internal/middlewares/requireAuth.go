package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/YashwanthN24/sell_point/internal/httpx"
	"github.com/golang-jwt/jwt/v5"
)

// middleware is a function that takes in a handler and returns a handler

const (
	authorization = "Authorization"
)

type AuthClaims struct {
	jwt.RegisteredClaims // embedding of struct composotion pattern
}

// STRUCT COMPOSITION
// https://chatgpt.com/c/6af431f1-c490-83ee-8185-9747fd6c33a5

func RequireAuth(logger *slog.Logger, secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.With("requestId", RequestIDFromContext(ctx))
			authHeader := r.Header.Get(authorization)

			if authHeader == "" {
				log.Error("no authorization header found")
				httpx.Error(w, http.StatusUnauthorized, "auth is required", httpx.CodeUnauthenticated)
				return
			}

			// tokenArr := strings.Split(authHeader, " ")
			// _, token := tokenArr[0], tokenArr[1]

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Error("no jwt token found ")
				httpx.Error(w, http.StatusUnauthorized, "auth is required", httpx.CodeUnauthenticated)
				return
			}

			token := parts[1]

			// verify token

			var claims AuthClaims
			_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")

				}
				return []byte(secret), nil
			}, jwt.WithValidMethods([]string{"HS256"}))

			if err != nil {
				log.Info("token rejected", "err", err)
				httpx.Error(w, http.StatusUnauthorized, "token is invalid or expired", httpx.CodeUnauthenticated)
				return
			}
			// pass the userid via context to tha handler
			// ctxWithUserId := context.WithValue(ctx , userIDKey , use)

			fmt.Println(claims)

			// TODO : Subject is userid (currently string but need to parse it to uuid before sending to the next handler)
			ctxWithUserId := context.WithValue(ctx, userIDKey, claims.Subject)

			// call the next handler
			next.ServeHTTP(w, r.WithContext(ctxWithUserId))

		})
	}
}

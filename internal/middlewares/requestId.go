package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// middleware is a function that takes in a handler and returns a handler

type ctxKey int

const (
	requestIDKey ctxKey = iota
)

const (
	requestId = "X-Request-ID"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestId)

		if id == "" {
			id = uuid.NewString()
		}

		// Add the request id backe to the response header
		w.Header().Add(requestId, id)

		// pass the request id via context to tha handler
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		// call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func RequestIDFromContext(ctx context.Context) string {
	requestId := ctx.Value(requestIDKey).(string)
	return requestId
}

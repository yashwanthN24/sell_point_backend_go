package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/YashwanthN24/sell_point/internal/httpx"
	"github.com/YashwanthN24/sell_point/internal/middlewares"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type MeHander struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewMeHandler(db *sql.DB, logger *slog.Logger) *MeHander {
	return &MeHander{
		db:     db,
		logger: logger,
	}
}

func (lh MeHander) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)
	// requestId := ctx.Value("requestCtxId").(string)
	userId, ok := middlewares.UserIDFromContext(ctx)
	if !ok {
		lh.logger.Error("user id not found in context", "requestId", requestId)
		httpx.Error(w, http.StatusInternalServerError, " something went wrong ", httpx.CodeInternalError)
		return
	}

	row := lh.db.QueryRowContext(ctx, `
			SELECT id , name , email , created_at  
			FROM users
			WHERE id = $1 
		`, userId)

	var user User

	err := row.Scan(&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpx.Error(w, http.StatusNotFound,
				"user not found", httpx.CodeNotFound)
			return
		}

		lh.logger.Error("failed to fetch user",
			"error", err,
			"requestId", requestId,
		)

		httpx.Error(w, http.StatusInternalServerError,
			"something went wrong", httpx.CodeInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

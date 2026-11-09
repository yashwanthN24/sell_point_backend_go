package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/YashwanthN24/sell_point/internal/dtos"
	"github.com/YashwanthN24/sell_point/internal/httpx"
	"github.com/YashwanthN24/sell_point/internal/middlewares"
	"github.com/jackc/pgx/v5/pgconn"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthHander struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewAuthHandler(db *sql.DB, logger *slog.Logger) *AuthHander {
	return &AuthHander{
		db:     db,
		logger: logger,
	}
}

func (ah AuthHander) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)

	logger := ah.logger.With("request_id", requestId)

	var req dtos.SignupRequest
	// Decoder to convert json to struct (desrializsation)

	// so now we are deserializing the payload (json body) into createListingRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode", "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid_body", httpx.CodeMalformedJSON)
		return
	}

	log.Println(req)

	if err := req.Validate(); err != nil {
		var verr *dtos.ValidationError
		errors.As(err, &verr)
		logger.Error("validation failed", "err", err)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	// has the password using bcrypt

	// cost is how slow you want your hashing function - 2^ 10 -> 1024 rounds of hashing (so if your db gets hacked also it will harder to try out different passwords as the algoritm ois slow)
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	if err != nil {
		logger.Error("hashing failed", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	row := ah.db.QueryRowContext(ctx,
		`INSERT INTO users (name , email , password) VALUES ($1 , $2 , $3) RETURNING id , created_at`, req.Name, req.Email, hash,
	)

	var u user
	err = row.Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		// based on error we can ay whether email dupliate dont need one more query to frst check if user exists and all
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			httpx.Error(w, http.StatusConflict, "email already taken", httpx.CodeConflict)
			return
		}

		logger.Error("scaning failed", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	out := dtos.SignupResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
	}

	logger.Info("new user registered", "user_id", out.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(out)
}

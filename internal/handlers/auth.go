package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/YashwanthN24/sell_point/internal/config"
	"github.com/YashwanthN24/sell_point/internal/dtos"
	"github.com/YashwanthN24/sell_point/internal/httpx"
	"github.com/YashwanthN24/sell_point/internal/middlewares"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var dummyHash = []byte("$2a$10$xOf7zVDwdjtj4qSFS4OfBeSF2W/pjnWF1nRKRxMbXEPiSFoKGu0G6")

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
	cfg    config.Config
}

func NewAuthHandler(db *sql.DB, logger *slog.Logger, cfg config.Config) *AuthHander {
	return &AuthHander{
		db:     db,
		logger: logger,
		cfg:    cfg,
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

func (ah AuthHander) Signin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)

	logger := ah.logger.With("request_id", requestId)

	var req dtos.SigninRequest
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

	row := ah.db.QueryRowContext(ctx,
		`SELECT id , email , password  FROM users 
		WHERE email LIKE $1`, req.Email, // can use = also in the sql query as we are finding exact match not a pattern
	)

	var user user
	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// to get same response time as success while brute forcig email password
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
			logger.Error("email not found try signing in first!!", "err", err)
			httpx.Error(w, http.StatusUnauthorized, "invalid credentials", httpx.CodeUnauthenticated)
			return
		}
		logger.Error("find user by email failed", "err", err)

		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	// check password match
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn("password mismatch", "user_id", user.ID)

		httpx.Error(w, http.StatusUnauthorized, "Invalid credentials", httpx.CodeUnauthenticated)
		return
	}

	// sign jwt token and return token to user

	tokenTTL := 24 * time.Hour
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   user.ID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(ah.cfg.JwtKey)) //secretky  muts be random strig maybe base64 and muts be very long in length

	if err != nil {
		logger.Error("jwt sign failed", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	out := dtos.SigninResponse{
		Token:     signedToken,
		ExpiresIn: int(tokenTTL.Seconds()),
	}

	logger.Info("new user logged in", "user_id", user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(out)
}

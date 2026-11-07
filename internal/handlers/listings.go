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
)

// add struct tags so we define how the fields are named when serialized to json using json.NewEncoder

// we do that by adding “ (backticks) `json:<name_of_field_you_want_in_json>`
type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

// constructor pattern to handle multiple dependencies for this handler at one place via struct and constructor oattern

type ListingHander struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHander {
	return &ListingHander{
		db:     db,
		logger: logger,
	}
}

// This pattern is called closure factory as inner func forms a closure with outer func
func (lh ListingHander) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)
	// requestId := ctx.Value("requestCtxId").(string)
	rows, err := lh.db.QueryContext(ctx, `
			SELECT id , title , description , price , created_at  
			FROM listings 
			ORDER BY created_at DESC 
			LIMIT 100
		`)

	if err != nil {
		lh.logger.Error("listings query error", "requestId", requestId, "err", err)
		// log.Printf("query: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	defer rows.Close() // to close the query connection (resource that the rows are holding eventually close them after this function ends )

	listings := []listing{}

	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.CreatedAt); err != nil {
			lh.logger.Error("rows scan error", "requestId", requestId, "err", err)
			// log.Printf("rows.scan : %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		lh.logger.Info("listings fetched", "total", len(listings), "requestId", requestId)

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows.err : %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listings) // this serializes our struct into json to be sent back to the client (Encode returns error which we are ignoring for the time being)

	// issue now it serialized exact as our struct filed exposing exact struct fields so we have to use struct fields to map them to how they look in json when serialized
}

func (lh ListingHander) Delete(w http.ResponseWriter, r *http.Request) {
	// request bound context always pass this forthrid part api and db calls as they avoid zombie queries i.e when client exits befoire the db finshes
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)
	// requestid := ctx.Value("requestCtxId").(string)
	id := r.PathValue("id")
	log.Println("id", id) // this also redirected to logger only as info because of our setting in main.go

	// db.Query() is only queries like select

	// for delete queries we use db.Exec()
	_, err := lh.db.ExecContext(ctx,
		`DELETE FROM listings WHERE id = $1`, id)

	if err != nil {
		// log.Printf("delete : %v", err)

		lh.logger.Error("delete failed", "listing_id", id, "requestId", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", "internal_error")
		// http.Error(w, "internal error", http.StatusInternalServerError) // return 500 internal server error as response
		return
	}

	// result.LastInsertId() helpful when  we have inserted data into table and want to get its id

	// rows , err :=result.RowsAffected()  // gives count of affected rows

	w.WriteHeader(http.StatusNoContent)
}

func (lh ListingHander) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)

	var req dtos.CreateListingRequest
	// Decoder to convert json to struct (desrializsation)

	// so now we are deserializing the payload (json body) into createListingRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("failed to decode", "requestId", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid_body", httpx.CodeMalformedJSON)
		return
	}

	// fmt.Println(req)
	log.Println(req)

	if err := req.Validate(); err != nil {
		var verr *dtos.ValidationError
		errors.As(err, &verr)
		lh.logger.Error("validation failed", "requestId", requestId, "err", err)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	row := lh.db.QueryRowContext(ctx, `
		INSERT INTO listings (title , description , price , city) 
		VALUES ($1 , $2 , $3 , $4) RETURNING id , title , created_at  `, req.Title, req.Description, req.Price, req.City)

	var out dtos.CreateListingResponse // because if we used listing struct then since all fields are not set we get zero values so elicitly use dto for both request snd response

	// to return all field modify to RETURNING * for all fields and take a varibale of type liting struct and pass their pointers to Scan
	if err := row.Scan(&out.ID, &out.Title, &out.CreatedAt); err != nil {
		lh.logger.Error("failed to insert", "requestId", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("listing created", "request_id", requestId, "listing_id", out.ID)

	// w.Write([]byte("created"))

	// response := struct {
	// 	Message string `json:"msg"`
	// }{Message: "Hello world"}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
	//  in go lang we dont have object so to send json we have to use structs as they get converted to json

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(out)
}

// https://chatgpt.com/c/6aef09fd-3a54-83ec-9112-5de682109dd4

package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"
)

// add struct tags so we define how the fields are named when serialized to json using json.NewEncoder

// we do that by adding “ (backticks) `json:<name_of_field_you_want_in_json>`
type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

// constructor pattern to handle multiple dependencies for this handler at one place via struct and constructor oattern

type ListingHander struct {
	db *sql.DB
}

func NewListingHandler(db *sql.DB) *ListingHander {
	return &ListingHander{
		db: db,
	}
}

// This pattern is called closure factory as inner func forms a closure with outer func
func (lh ListingHander) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx, `
			SELECT id , title , description , price , created_at  
			FROM listings 
			ORDER BY created_at DESC 
			LIMIT 100
		`)

	if err != nil {
		log.Printf("query: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	defer rows.Close() // to close the query connection (resource that the rows are holding eventually close them after this function ends )

	listings := []listing{}

	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.CreatedAt); err != nil {
			log.Printf("rows.scan : %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

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
	id := r.PathValue("id")
	log.Println("id", id)

	// db.Query() is only queries like select

	// for delete queries we use db.Exec()
	_, err := lh.db.ExecContext(ctx,
		`DELETE FROM listings WHERE id = $1`, id)

	if err != nil {
		// log.Printf("delete : %v", err)

		// using global logger setup across app bad practice always we must be clear in which hander needs what dependency like how we do for db  so thats //TODO
		slog.Error("delete failed", "listing_id", id, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError) // return 500 internal server error as response
		return
	}

	// result.LastInsertId() helpful when  we have inserted data into table and want to get its id

	// rows , err :=result.RowsAffected()  // gives count of affected rows

	w.WriteHeader(http.StatusNoContent)
}

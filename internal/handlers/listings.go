package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// add struct tags so we define how the fields are named when serialized to json using json.NewEncoder 

// we do that by adding `` (backticks) `json:<name_of_field_you_want_in_json>`  
type listing struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
	City        string `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}


// This pattern is called closure factory as inner func forms a closure with outer func
func List(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query(`
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
				log.Printf("rows.scan : %w", err)
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

		w.Header().Set("Content-Type" , "application/json")
		w.WriteHeader(http.StatusOK)


		_ = json.NewEncoder(w).Encode(listings) // this serializes our struct into json to be sent back to the client (Encode returns error which we are ignoring for the time being)

		// issue now it serialized exact as our struct filed exposing exact struct fields so we have to use struct fields to map them to how they look in json when serialized
	}
}

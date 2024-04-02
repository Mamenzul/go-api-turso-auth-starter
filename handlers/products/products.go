package products

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/nrednav/cuid2"

	models "github.com/mamenzul/go-api/models"
)

func Router(db *sql.DB) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := cuid2.Generate()
		result, err := db.Exec("INSERT INTO products (id, code, price) VALUES (?, ?, ?)", id, product.Code, product.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Product created successfully",
			"id":           id,
			"rowsAffected": rowsAffected,
		})

	})

	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		var products []models.Product
		rows, err := db.Query("SELECT code, price FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var product models.Product
			err := rows.Scan(&product.Code, &product.Price)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(products)

	})
	return r
}

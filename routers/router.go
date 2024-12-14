package router

import (
	"net/http"
	"github.com/BlackDevil559/novahack2/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.Use(corsMiddleware)
	router.HandleFunc("/api/addnewuser", controllers.AddNewUser).Methods("POST")
	router.HandleFunc("/api/addnewfood", controllers.AddNewFood).Methods("POST")
	router.HandleFunc("/api/deleteuser/{id}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/deletefood/{id}", controllers.DeleteFoodItem).Methods("DELETE")
	router.HandleFunc("/api/showfood/{id}", controllers.ShowFoodNearBy).Methods("GET")
	router.HandleFunc("/api/bookfood/{id}/{ns}/{consumerid}", controllers.BookFooditem).Methods("GET")
	router.HandleFunc("/api/login/{phone}/{password}", controllers.LoginUser).Methods("GET","")
	router.HandleFunc("/api/showallorder/{id}", controllers.ShowAllOrder).Methods("GET","OPTIONS")
	router.HandleFunc("/api/addrating/{id}/{rating}", controllers.AddRating).Methods("PUT","OPTIONS")
	router.HandleFunc("/api/showallfood", controllers.ShowAllFood).Methods("GET", "OPTIONS")
	return router
}

// CORS middleware to handle OPTIONS and other headers for CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow any origin, but you can specify a specific domain instead of "*"
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("User-Agent","CustomUserAgent/1.0")
		// If it's a preflight OPTIONS request, return a status 200 OK
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed with the actual request
		next.ServeHTTP(w, r)
	})
}

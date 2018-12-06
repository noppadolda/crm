package routers

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	generic "ats_eng_api/generic"
	atsHandler "ats_eng_api/handler"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitRouter() *mux.Router {
	/* allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}) */

	r := mux.NewRouter()
	r.HandleFunc("/atsapi/{category}/{function}", atsHandler.HandleRequest).Methods("GET", "POST", "PUT", "OPTIONS")
	r.HandleFunc("/atsapi/{function}", atsHandler.HandleRequest).Methods("GET", "POST", "PUT", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*", "*:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept, Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	rs := c.Handler(r)

	/* srv := &http.Server{
		Handler:      handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(r),
		Addr:         ":" + generic.GetAPIGenericValue("port"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	} */
	srv := &http.Server{
		Handler:      rs,
		Addr:         ":" + generic.GetAPIGenericValue("port"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//fmt.Println("start port 2918...")
	log.Println("start server " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
	return r
}

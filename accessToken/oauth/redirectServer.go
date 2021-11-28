package oauth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type redirectServer struct {
	oAuthState   string
	redirectUri  string
	redirectPort int
}

func newRedirectServer(oAuthState string, redirectUri string, redirectPort int) *redirectServer {
	return &redirectServer{
		oAuthState:   oAuthState,
		redirectUri:  redirectUri,
		redirectPort: redirectPort,
	}
}

func (s *redirectServer) GetCode() string {
	log.Info("Starting redirecting server")
	code := ""

	router := mux.NewRouter()

	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", s.redirectPort),
	}

	router.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(301)
		w.Header().Add("Location", s.redirectUri)

		stateQuery := r.URL.Query().Get("state")
		codeQuery := r.URL.Query().Get("code")

		if stateQuery != s.oAuthState {
			// shutdown server
			server.Shutdown(context.Background())
			log.Fatal("received OAuth state does not match the send state")
		}

		code = codeQuery

		// shutdown server
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("unable to shutdown server. Error: %s\n", err.Error())
		}
	})

	log.Info("Waiting for you to grant access...")

	// this will block until we get a response from pinterest
	server.ListenAndServe()

	return code

}

package passwordless

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

// BuildRoutes returns the routes for the application.
func BuildRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/login-success", LoginSuccessHandler)
	router.HandleFunc("/verify", VerifyHandler)
	router.HandleFunc("/logout", LogoutHandler)

	// profile routes with LoginRequiredMiddleware
	profileRouter := mux.NewRouter()
	profileRouter.HandleFunc("/profile", ProfileHandler)

	router.PathPrefix("/profile").Handler(negroni.New(
		negroni.HandlerFunc(LoginRequiredMiddleware),
		negroni.Wrap(profileRouter),
	))

	// apply the base middleware to the main router
	n := negroni.New(
		negroni.HandlerFunc(CsrfMiddleware),
		negroni.HandlerFunc(UserMiddleware),
	)
	n.UseHandler(router)

	return n
}

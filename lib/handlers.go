package passwordless

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
)

// loginEmailTemplate is used in verification emails.
var loginEmailTemplate = template.Must(template.New("loginEmailTemplate").Parse(`
  A login request was submitted for this email in the passwordless demo app. Use the link below to verify:

  {{ .LoginUrl }}`))

// redirectToOrigin redirects the user to their OriginUrl if set.
func redirectToOrigin(user *User, w http.ResponseWriter, r *http.Request) {
	var redirectUrl string
	if user.OriginUrl.String != "" {
		redirectUrl = user.OriginUrl.String
	} else {
		redirectUrl = "/profile"
	}
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect logged in users
	user := GetContextUser(r)
	if user != nil {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		user := &User{}

		doRedirect := func() {
			http.Redirect(w, r, "/login-success", http.StatusFound)
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		email := r.PostFormValue("Email")
		if email != "" {
			// Get or create user with given email
			if err := dbmap.SelectOne(user, "SELECT * FROM users WHERE email=$1", email); err != nil {
				user.Email = email
				user.RefreshToken()

				if err := dbmap.Insert(user); err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"email": email,
					}).Warn("Error creating user.")
					doRedirect()
					return
				}
			} else {
				// Update existing user's token
				user.RefreshToken()
				dbmap.Update(user)
			}

			// Update user.OriginUrl with referring page if its from this host, assuming
			// they came here via a redirect trying to access a page that requires auth
			if referrerUrl, err := url.ParseRequestURI(r.Referer()); err != nil {
				if referrerUrl.Scheme == r.URL.Scheme && referrerUrl.Host == r.URL.Host {
					if err := user.UpdateOriginUrl(referrerUrl); err != nil {
						dbmap.Update(user)
					}
				}
			}

			// Build login url
			params := url.Values{}
			params.Add("token", user.Token)
			params.Add("uid", strconv.FormatInt(user.Id, 10))

			loginUrl := url.URL{}

			if r.URL.IsAbs() {
				loginUrl.Scheme = r.URL.Scheme
				loginUrl.Host = r.URL.Host
			} else {
				loginUrl.Scheme = "http"
				loginUrl.Host = r.Host
			}

			loginUrl.Path = "/verify"

			// Send login email
			var mailContent bytes.Buffer
			ctx := struct {
				LoginUrl string
			}{
				fmt.Sprintf("%s?%s", loginUrl.String(), params.Encode()),
			}

			go func() {
				if err := loginEmailTemplate.Execute(&mailContent, ctx); err == nil {
					if err := SendMail([]string{email}, "Passwordless Login Verification", mailContent.String()); err != nil {
						log.WithFields(log.Fields{
							"error": err,
						}).Error("Error sending verification email")
					}
				}
			}()
		}

		doRedirect()
		return
	}

	renderTemplate(w, r, "home", nil)
}

func LoginSuccessHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect logged in users
	user := GetContextUser(r)
	if user != nil {
		redirectToOrigin(user, w, r)
		return
	}

	renderTemplate(w, r, "login-success", nil)
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect logged in users
	user := GetContextUser(r)
	if user != nil {
		redirectToOrigin(user, w, r)
		return
	}

	// Collect URL params
	params := r.URL.Query()
	userId := params.Get("uid")
	userToken := params.Get("token")

	doResponse := func() {
		// Something failed along the way...
		renderTemplate(w, r, "verify", nil)
	}

	if userId != "" && userToken != "" {
		userId, err := strconv.ParseInt(userId, 0, 64)
		if err != nil {
			doResponse()
			return
		}

		if obj, err := dbmap.Get(User{}, userId); err == nil {
			user := obj.(*User)
			if user.IsValidToken(userToken) {
				// Valid token, log user in
				Login(user, w, r)

				// Do redirect
				redirectToOrigin(user, w, r)
				return
			}
		}
	}

	doResponse()
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	Logout(w, r)
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "profile", nil)
}

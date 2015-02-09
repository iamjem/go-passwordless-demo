package passwordless

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var templates *template.Template

type tmplContext struct {
	CsrfToken string
	Request   *http.Request
	Page      map[string]interface{}
}

func newTmplContext(r *http.Request, d map[string]interface{}) *tmplContext {
	ctx := tmplContext{
		context.Get(r, "csrf_token").(string),
		r,
		d,
	}
	return &ctx
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, d map[string]interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, newTmplContext(r, d))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	// load all templates
	root, _ := os.Getwd()

	tmplNames := []string{}

	err := filepath.Walk(path.Join(root, "lib", "templates"), func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			tmplNames = append(tmplNames, path)
		}
		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Template load error.")
	}

	templates = template.New("")
	if _, err := templates.ParseFiles(tmplNames...); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Template parsing error.")
	}

}

package main
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/alexedwards/scs/v2"
	"insta/pkg/config"
	"insta/pkg/handler"
	"insta/pkg/render"
)
var app config.AppConfig
var session *scs.SessionManager
func main() {
	app.InProduction = true
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
	}
	app.TemplateCache = tc
	app.UseCache = true
	repo := handler.NewRepo(&app)
	handler.NewHandlers(repo)
	render.NewTemplates(&app)
	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		portNumber = "8080"
	}
	fmt.Printf("Starting production application on port %s\n", portNumber)
	srv := &http.Server{
		Addr:    ":" + portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

One more important check

Make sure your route is uncommented:

mux.HandleFunc("/", handler.Repo.Index)

And your Index is:

func (m *Repository) Index(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "index.html")
}

After committing this to GitHub, Render should automatically rebuild if automatic deploys are enabled.

If Safari still downloads an empty file, don’t change anything else yet—send me the new Render log from startup.

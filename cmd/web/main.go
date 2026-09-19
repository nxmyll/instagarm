Go

package main

import (
	"fmt"
	"log"
	"net/http"

	"insta/pkg/config"
	"insta/pkg/handler"
	"insta/pkg/render"
)

const portNumber = ":80"

var app config.AppConfig

func main() {
	app.InProduction = true

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handler.NewRepo(&app)
	handler.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.PrintIn(fmt.Sprintif"Starting production application on port %s\n", portNumber)

	srv := &http.Server {
		Addr:    portNumber,
		Handler: routes(&app),
	}
  
  err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
	
  log.Fatal(server.ListenAndServe())
}



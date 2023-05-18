package router

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/training/pkg/model"
	"github.com/training/pkg/param"
	"github.com/training/pkg/render"
)

func HandlerHTTP(r *chi.Mux) {
	r.Use(noSurf)
	r.Use(sessionLoad)
	r.Use(middleware.Logger)
	r.Get("/", home)
	r.Get("/about", about)
}

func home(w http.ResponseWriter, r *http.Request) {
	p := param.Eject(r)
	remoteIP := r.RemoteAddr
	s := p.Session
	s.Put(r.Context(), "remote_ip", remoteIP)
	out, err := render.Client("home.page.tmpl", &model.TemplateData{}, true).RenderTemplate()
	if err != nil {
		return
	}
	defer out.Close()
	io.Copy(w, out)
}

func about(w http.ResponseWriter, r *http.Request) {
	p := param.Eject(r)
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"
	remoteIP := p.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	out, err := render.Client("about.page.tmpl", &model.TemplateData{
		StringMap: stringMap,
	}, true).RenderTemplate()
	if err != nil {
		return
	}
	defer out.Close()
	io.Copy(w, out)
}

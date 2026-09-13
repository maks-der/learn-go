package web

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"learngo/internal/content"
)

const langCookie = "learngo-lang"

// Server serves handbook pages from Markdown files.
type Server struct {
	store *content.Store
	tmpl  *template.Template
	mux   *http.ServeMux
}

type page struct {
	Lang      string
	OtherLang string
	HomeURL   string
	SelfURL   string
	AltURL    string
	UI        UI
	Topics    []content.Topic
	Topic     *content.Topic
	Prev      *content.Topic
	Next      *content.Topic
	Status    int
}

// New builds the HTTP handler.
func New(store *content.Store, root string) (http.Handler, error) {
	tmpl, err := template.ParseGlob(filepath.Join(root, "web", "templates", "*.html"))
	if err != nil {
		return nil, err
	}
	s := &Server{store: store, tmpl: tmpl, mux: http.NewServeMux()}
	static := http.FileServer(http.Dir(filepath.Join(root, "web", "static")))
	s.mux.Handle("GET /static/", http.StripPrefix("/static/", static))
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /{$}", s.redirectHome)
	s.mux.HandleFunc("GET /en/{$}", s.home)
	s.mux.HandleFunc("GET /ru/{$}", s.home)
	s.mux.HandleFunc("GET /en/{slug}", s.topic)
	s.mux.HandleFunc("GET /ru/{slug}", s.topic)
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

func (s *Server) redirectHome(w http.ResponseWriter, r *http.Request) {
	lang := "en"
	if c, err := r.Cookie(langCookie); err == nil && content.ValidLang(c.Value) {
		lang = c.Value
	}
	http.Redirect(w, r, "/"+lang+"/", http.StatusSeeOther)
}

func langFromPath(r *http.Request) string {
	path := strings.Trim(r.URL.Path, "/")
	if i := strings.IndexByte(path, '/'); i >= 0 {
		path = path[:i]
	}
	return strings.ToLower(path)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	lang := langFromPath(r)
	if !content.ValidLang(lang) {
		s.notFound(w, r, "en")
		return
	}
	setLangCookie(w, lang)
	p := page{
		Lang:      lang,
		OtherLang: otherLang(lang),
		HomeURL:   "/" + lang + "/",
		SelfURL:   "/" + lang + "/",
		AltURL:    "/" + otherLang(lang) + "/",
		UI:        uiFor(lang),
		Topics:    s.store.List(lang),
	}
	s.render(w, "home", p, http.StatusOK)
}

func (s *Server) topic(w http.ResponseWriter, r *http.Request) {
	lang := langFromPath(r)
	slug := r.PathValue("slug")
	if !content.ValidLang(lang) || !content.ValidSlug(slug) {
		s.notFound(w, r, fallbackLang(lang))
		return
	}
	topic, ok := s.store.Get(lang, slug)
	if !ok {
		s.notFound(w, r, lang)
		return
	}
	setLangCookie(w, lang)
	prev, next := s.store.Neighbors(lang, slug)
	t := topic
	p := page{
		Lang:      lang,
		OtherLang: otherLang(lang),
		HomeURL:   "/" + lang + "/",
		SelfURL:   "/" + lang + "/" + slug,
		AltURL:    "/" + otherLang(lang) + "/" + slug,
		UI:        uiFor(lang),
		Topic:     &t,
		Prev:      prev,
		Next:      next,
	}
	if _, exists := s.store.Get(p.OtherLang, slug); !exists {
		p.AltURL = "/" + p.OtherLang + "/"
	}
	s.render(w, "topic", p, http.StatusOK)
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request, lang string) {
	if !content.ValidLang(lang) {
		lang = "en"
	}
	p := page{
		Lang:      lang,
		OtherLang: otherLang(lang),
		HomeURL:   "/" + lang + "/",
		SelfURL:   r.URL.Path,
		AltURL:    "/" + otherLang(lang) + "/",
		UI:        uiFor(lang),
		Status:    http.StatusNotFound,
	}
	s.render(w, "notfound", p, http.StatusNotFound)
}

func (s *Server) render(w http.ResponseWriter, name string, p page, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	if err := s.tmpl.ExecuteTemplate(w, name, p); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func otherLang(lang string) string {
	if lang == "ru" {
		return "en"
	}
	return "ru"
}

func fallbackLang(lang string) string {
	if content.ValidLang(lang) {
		return lang
	}
	return "en"
}

func setLangCookie(w http.ResponseWriter, lang string) {
	http.SetCookie(w, &http.Cookie{
		Name:     langCookie,
		Value:    lang,
		Path:     "/",
		MaxAge:   365 * 24 * 3600,
		HttpOnly: false,
		Secure:   os.Getenv("RENDER") != "",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})
}

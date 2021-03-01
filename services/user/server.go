package user

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iconmobile-dev/backend-coding-challenge/config"
	"github.com/iconmobile-dev/backend-coding-challenge/lib/bootstrap"
	"github.com/iconmobile-dev/backend-coding-challenge/lib/handlers"
	"github.com/iconmobile-dev/backend-coding-challenge/lib/storage"
	"github.com/iconmobile-dev/go-core/logger"
)

var log logger.Logger
var cfg config.Config

// SetupLoggerAndConfig sets the global logger and config dependency
// should be called during tests
func SetupLoggerAndConfig(serverName string, test bool) {
	log, cfg = bootstrap.LoggerAndConfig(serverName, test)
	storage.SetupLoggerAndConfig(serverName, test)
}

// initiates log and cfg with default values
func init() {
	SetupLoggerAndConfig("engagement", false)
}

// Server manages the internal state of the service.
type Server struct {
	db     *storage.DB
	cache  *storage.Cache
	router *chi.Mux
}

// New provisions the service defaults: storage database, cache, routes
func New(db *storage.DB, cache *storage.Cache) *Server {
	r := chi.NewRouter()
	handlers.DefaultMiddlewares(r)

	s := &Server{
		db:     db,
		cache:  cache,
		router: r,
	}

	// load routes
	s.routes()

	return s
}

// ServeHTTP makes our server a http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

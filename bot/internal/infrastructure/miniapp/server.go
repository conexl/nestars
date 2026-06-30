package miniapp

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/fx"

	"nestars-bot/internal/config"
	"nestars-bot/internal/domain"
	"nestars-bot/internal/usecase"
)

// Server serves the Mini App static files and API.
type Server struct {
	cfg       *config.Config
	userUC    *usecase.User
	catalog   *usecase.Catalog
	order     *usecase.Order
	orderRepo domain.OrderRepository
}

// NewServer creates a new Mini App server.
func NewServer(
	cfg *config.Config,
	userUC *usecase.User,
	catalog *usecase.Catalog,
	order *usecase.Order,
	orderRepo domain.OrderRepository,
) *Server {
	return &Server{
		cfg:       cfg,
		userUC:    userUC,
		catalog:   catalog,
		order:     order,
		orderRepo: orderRepo,
	}
}

// Module registers the Mini App HTTP server in the fx application.
var Module = fx.Module("miniapp",
	fx.Provide(NewServer),
	fx.Invoke(Run),
)

func (s *Server) router() http.Handler {
	mux := http.NewServeMux()

	// API routes.
	mux.HandleFunc("POST /api/miniapp/auth", s.handleAuth)
	mux.HandleFunc("GET /api/miniapp/catalog/stars", s.authMiddleware(s.handleCatalogStars))
	mux.HandleFunc("GET /api/miniapp/catalog/premium", s.authMiddleware(s.handleCatalogPremium))
	mux.HandleFunc("POST /api/miniapp/orders", s.authMiddleware(s.handleCreateOrder))
	mux.HandleFunc("POST /api/miniapp/orders/", s.authMiddleware(s.handlePayOrder))
	mux.HandleFunc("GET /api/miniapp/profile", s.authMiddleware(s.handleProfile))
	mux.HandleFunc("GET /api/miniapp/orders", s.authMiddleware(s.handleOrders))
	mux.HandleFunc("POST /api/miniapp/withdraw", s.authMiddleware(s.handleWithdraw))

	// Static files (SPA).
	staticPath := s.cfg.MiniAppStaticPath
	if staticPath == "" {
		staticPath = "./miniapp/dist"
	}
	mux.HandleFunc("/", s.spaHandler(staticPath))

	return mux
}

// Run starts the Mini App HTTP server as an fx lifecycle hook.
func Run(lc fx.Lifecycle, s *Server) {
	server := &http.Server{
		Addr:    ":" + s.cfg.MiniAppPort,
		Handler: s.router(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("mini app server listening on %s", server.Addr)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("mini app server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func (s *Server) spaHandler(staticPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || filepath.Ext(r.URL.Path) == "" {
			http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
			return
		}

		path := filepath.Join(staticPath, r.URL.Path)
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
			return
		}

		http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r)
	}
}

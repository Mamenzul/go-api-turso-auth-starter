package api

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	authService "github.com/mamenzul/go-api/services/auth"
	"github.com/mamenzul/go-api/services/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {

	userStore := user.NewStore(s.db)

	options := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) { // secret key for JWT
			return "secret", nil
		}),
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		Issuer:         "go-api",
		URL:            "http://127.0.0.1:8080",
		AvatarStore:    avatar.NewLocalFS("/tmp"),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
			// allow only dev_* names
			return claims.User != nil && strings.HasPrefix(claims.User.Name, "dev_")
		}),
	}

	// create auth service with providers
	service := auth.NewService(options)
	service.AddDirectProvider("local", provider.CredCheckerFunc(func(user, password string) (ok bool, err error) {
		u, err := userStore.GetUserByEmail(user)
		if err != nil {
			return false, err
		}
		ok = authService.ComparePasswords(u.Password, []byte(password))
		return ok, err
	}))

	// setup http server
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// setup auth routes
	authRoutes, avaRoutes := service.Handlers()
	router.Mount("/auth", authRoutes) // add auth handlers
	router.Mount("/avatar", avaRoutes)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router) // add avatar handler
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	return http.ListenAndServe(":8080", router)
}

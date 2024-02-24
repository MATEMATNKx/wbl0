package transport

import (
	"context"
	"l0/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type OrderService interface {
	Create(orderUID, data string)
	Get(orderUID string) (string, error)
}

type StanService interface {
	Start()
	Stop()
}

type App struct {
	cfg *config.Config

	web      *http.Server
	router   *gin.Engine
	orderSvc OrderService
	stanSvc  StanService
}

func New(cfg *config.Config, orderSvc OrderService, stanSvc StanService) *App {

	router := gin.New()

	app := &App{
		cfg:      cfg,
		router:   router,
		orderSvc: orderSvc,
		stanSvc:  stanSvc,
	}

	app.web = &http.Server{
		Addr:    app.cfg.Web.Address(),
		Handler: router,
	}

	app.initRoutes()

	return app
}

func (app *App) initRoutes() {
	app.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Context-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	app.router.GET("/:id", app.GetById)

}
func (app *App) Start() error {
	app.stanSvc.Start()
	log.Printf("start web server at http://%s\n", app.web.Addr)
	return app.web.ListenAndServe()
}
func (app *App) Stop(ctx context.Context) error {
	app.stanSvc.Stop()
	log.Println("stop web server")
	return app.web.Shutdown(ctx)
}

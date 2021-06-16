package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rtbe/clean-rest-api/delivery/web/handlers"
	mid "github.com/rtbe/clean-rest-api/delivery/web/middlewares"
	"github.com/rtbe/clean-rest-api/domain/usecase"
	"github.com/rtbe/clean-rest-api/internal/logger"
)

// App handles interaction between web router and business case services.
type App struct {
	router   *chi.Mux
	services usecase.Services
	logger   logger.Logger
}

// NewApp creates a new application.
func NewApp(s usecase.Services, l logger.Logger) *App {

	r := chi.NewMux()

	// Set up middlewares for whole application:
	r.Use(mid.RequestInfo, mid.Logger(l))

	// Configure routes for Auth Group
	ag := handlers.AuthGroup{AuthService: s.Auth}
	r.With().Route("/auth", func(r chi.Router) {
		r.Method(http.MethodPost, "/signup", handlers.Handler{H: ag.SignUp, L: l})
		r.With().Method(http.MethodPost, "/signin", handlers.Handler{H: ag.SignIn, L: l})
		r.With().Method(http.MethodPost, "/signout", handlers.Handler{H: ag.SignOut, L: l})
		r.With().Method(http.MethodPost, "/refresh", handlers.Handler{H: ag.RefreshTokens, L: l})
	})

	// Configure routes for User Group
	ug := handlers.UserGroup{UserService: s.User}
	r.With().Route("/users", func(r chi.Router) {
		r.With().Method(http.MethodGet, "/{lastSeenID}/{limit}", handlers.Handler{H: ug.ListUsers, L: l})
		r.Method(http.MethodGet, "/{id}", handlers.Handler{H: ug.GetUserByID, L: l})
		r.With().Method(http.MethodPatch, "/{id}", handlers.Handler{H: ug.UpdateUser, L: l})
		r.With().Method(http.MethodDelete, "/{id}", handlers.Handler{H: ug.DeleteUser, L: l})
	})

	// Configure routes for Product Group
	pg := handlers.ProductGroup{ProductService: s.Product}
	r.With().Route("/products", func(r chi.Router) {
		r.With().Method(http.MethodPost, "/", handlers.Handler{H: pg.CreateProduct, L: l})
		r.Method(http.MethodGet, "/{lastSeenID}/{limit}", handlers.Handler{H: pg.ListProducts, L: l})
		r.Method(http.MethodGet, "/{id}", handlers.Handler{H: pg.GetProduct, L: l})
		r.With().Method(http.MethodPatch, "/{id}", handlers.Handler{H: pg.UpdateProduct, L: l})
		r.Method(http.MethodDelete, "/{id}", handlers.Handler{H: pg.DeleteProduct, L: l})
	})

	// Configure routes for Order Group
	og := handlers.OrderGroup{OrderService: s.Order}
	r.With().Route("/orders", func(r chi.Router) {
		r.With().Method(http.MethodPost, "/", handlers.Handler{H: og.CreateOrder, L: l})
		r.Method(http.MethodGet, "/{id}", handlers.Handler{H: og.GetOrder, L: l})
		r.Method(http.MethodGet, "/{lastSeenID}/{limit}", handlers.Handler{H: og.ListOrders, L: l})
		r.With().Method(http.MethodPatch, "/{id}", handlers.Handler{H: og.UpdateOrder, L: l})
		r.Method(http.MethodDelete, "/{id}", handlers.Handler{H: og.DeleteOrder, L: l})
		r.Route("/users", func(r chi.Router) {
			r.Method(http.MethodGet, "/{userID}", handlers.Handler{H: og.ListUserOrders, L: l})
			r.Method(http.MethodDelete, "/{userID}", handlers.Handler{H: og.DeleteUserOrders, L: l})
		})
	})

	// Configure routes for Order Items Group
	oig := handlers.OrderItemGroup{OrderItemService: s.OrderItem}
	r.With().Route("/order_items", func(r chi.Router) {
		r.With().Method(http.MethodPost, "/", handlers.Handler{H: oig.CreateOrderItem, L: l})
		r.Method(http.MethodGet, "/{id}", handlers.Handler{H: oig.GetOrderItem, L: l})
		r.With().Method(http.MethodPatch, "/{id}", handlers.Handler{H: oig.UpdateOrderItem, L: l})
		r.Method(http.MethodDelete, "/{id}", handlers.Handler{H: oig.DeleteOrderItem, L: l})
		// TODO: WHERE SHOULD I PUT EM?
		r.Route("/orders/{orderID}", func(r chi.Router) {
			r.Method(http.MethodGet, "/", handlers.Handler{H: oig.ListOrderOrderItems, L: l})
			r.Method(http.MethodDelete, "/", handlers.Handler{H: oig.DeleteOrderOrderItems, L: l})
		})
	})

	// Configure routes for Status Group
	sg := handlers.StatusGroup{}
	r.Method(http.MethodGet, "/status", handlers.Handler{H: sg.Status, L: l})

	// Configure routes for Documentation
	handlerSwagger := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.yaml")
	}
	r.HandleFunc("/swagger.yaml", handlerSwagger)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	redocHandler := middleware.Redoc(opts, nil)

	// Documentation will be available here
	r.Handle("/docs", redocHandler)

	app := App{
		router:   r,
		services: s,
		logger:   l,
	}

	return &app
}

// ServeHTTP lets an App implements http.Handler interface.
// So it is an entry point to all incoming http traffic.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

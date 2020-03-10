package app

import "net/http"

// за разделение handler'ов по адресам -> routing
func (receiver *server) InitRoutes() {
	mux := receiver.router.(*exactMux)
	// panic, если происходят конфликты
	// Handle - добавляет Handler (неудобно)
	// HandleFunc

	// стандартный mux:
	// - если адрес начинается со "/" - то под действие обработчика попадает всё, что начинается со "/"
	// https://dropmefiles.com/k0P8d
	mux.GET("/", receiver.handleBurgersList(false))
	mux.POST("/", receiver.handleBurgersList(false))

	mux.GET("/admin/burgers", receiver.handleBurgersList(true))
	mux.POST("/admin/burgers/save", receiver.handleBurgersSave())
	mux.POST("/admin/burgers/remove", receiver.handleBurgersRemove())

	mux.GET("/slow", receiver.handleSlow())

	// - но если есть более "специфичный", то используется он
	mux.GET("/favicon.ico", receiver.handleFavicon())
	mux.GET("/media", http.StripPrefix("/media", http.FileServer(http.Dir(receiver.mediaPath))).ServeHTTP)
}

package composition

// Application initialization

func Run() {
	config_init() // TODO Initialize configuration
	email_init()  // TODO Initialize email
	db_init()     // TODO Initialize database
	http_init()   // TODO Initialize HTTP server

	select {}
}

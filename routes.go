package main

func (s *server) routes() {
	v1 := s.router.Group("/v1")

	v1.GET("/", s.handlerHome)
	v1.GET("/home", s.handlerHome)
	v1.GET("/json", s.handlerBigJSON)
	v1.GET("/json-stream", s.handlerBigJSONStream)

	// Users routes
	// ------------
	users := v1.Group("/users")
	users.GET("/", s.handlerGetUser)
}

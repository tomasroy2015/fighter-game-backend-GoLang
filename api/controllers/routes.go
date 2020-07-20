package controllers

import "gitlab.com/zenport.io/go-assignment/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//knight routes
	s.Router.HandleFunc("/knight", middlewares.SetMiddlewareJSON(s.CreateKnight)).Methods("POST")
	s.Router.HandleFunc("/knight", middlewares.SetMiddlewareJSON(s.GetKnights)).Methods("GET")
	s.Router.HandleFunc("/knight/{id}", middlewares.SetMiddlewareJSON(s.GetKnight)).Methods("GET")
}

package controllers

import (
	"net/http"

	"gitlab.com/zenport.io/go-assignment/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To fighter  API")

}

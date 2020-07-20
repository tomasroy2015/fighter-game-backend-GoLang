package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/zenport.io/go-assignment/api/models"
	"gitlab.com/zenport.io/go-assignment/api/responses"
)

func (server *Server) CreateKnight(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	knight := models.Knight{}
	err = json.Unmarshal(body, &knight)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	knight.Prepare()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	knightCreated, err := knight.SaveKnight(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, knightCreated.ID))
	responses.JSON(w, http.StatusCreated, knightCreated)
}

func (server *Server) GetKnights(w http.ResponseWriter, r *http.Request) {

	knight := models.Knight{}

	knights, err := knight.FindAllKnights(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, knights)
}

func (server *Server) GetKnight(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	knight := models.Knight{}

	knightReceived, err := knight.FindKnightByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, knightReceived)
}

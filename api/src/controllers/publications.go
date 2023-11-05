package controllers

import (
	"api/src/auth"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePublication(w http.ResponseWriter, r *http.Request) {
	userId, erro := auth.ExtractUserId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	var publication models.Publication
	if erro = json.NewDecoder(r.Body).Decode(&publication); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	publication.AuthorId = userId

	if erro = publication.Prepare(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositoriy := repositories.NewPublicationRepository(db)
	newPublication, erro := repositoriy.Create(publication)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, newPublication)
}

func GetPublication(w http.ResponseWriter, r *http.Request) {

}

func GetPublicationById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	publicationId, erro := strconv.ParseUint(params["publicationId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repository := repositories.NewPublicationRepository(db)
	publication, erro := repository.GetPublicationById(publicationId)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publication)
}

func UpdatePublication(w http.ResponseWriter, r *http.Request) {

}

func DeletePublication(w http.ResponseWriter, r *http.Request) {

}



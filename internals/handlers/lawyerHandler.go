package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Nathac/go-api/internals/store"
	"github.com/Nathac/go-api/internals/utils"
)

type LawyerHandler struct {
	LawyerStore store.LawyerInterface
	Logger      *log.Logger
}

func NewLawyerHandler(LawyerStore store.LawyerInterface, logger *log.Logger) *LawyerHandler {
	handler := &LawyerHandler{
		LawyerStore: LawyerStore,
		Logger:      logger,
	}
	return handler
}

func (lw *LawyerHandler) CreateLawyer(w http.ResponseWriter, r *http.Request) {

	var Lawyer store.Lawyer
	err := json.NewDecoder(r.Body).Decode(&Lawyer)

	if err != nil {
		lw.Logger.Printf("error decoding the lawyer: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"errors": "error in decoding the lawyer"})
		return
	}

	CreatedLawyer, err := lw.LawyerStore.CreateLawyer(&Lawyer)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"errors": "error executing to db"})
		return
	}
	utils.WriteJson(w, http.StatusOK, utils.Envelop{"lawyer": CreatedLawyer})
}

func (lw *LawyerHandler) GetLawyerById(w http.ResponseWriter, r *http.Request) {

	Id, err := utils.IdReader(r)
	if err != nil {
		lw.Logger.Printf("ERROR READING ID: %v", err)
		utils.WriteJson(w, http.StatusNoContent, utils.Envelop{"errors": "ERROR READING ID"})
		return
	}

	existingLawyer, err := lw.LawyerStore.GetLawyerById(Id)
	if err != nil {
		lw.Logger.Printf("ERROR PERFORMING OPERATION WITH THE ID: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"errors": "ERROR PERFORMING OPERATION:"})
		return
	}

	lw.Logger.Printf("successfully fetched lawyer with the id of  %v", Id)
	utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"lawyer": existingLawyer})

}

//
//
//

func (lw *LawyerHandler) UpdateLawyer(w http.ResponseWriter, r *http.Request) {
	var lawyer store.Lawyer
	err := json.NewDecoder(r.Body).Decode(&lawyer)
	if err != nil {
		lw.Logger.Printf("error in decoding the response: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"errors": "decoding the data:"})
		return
	}

	err = lw.LawyerStore.UpdateLawyer(&lawyer)
	if err != nil {
		lw.Logger.Printf("error updating...: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"errors": "error updating:"})
		return
	}

	fmt.Fprint(w, "lawyer updated")
}

//
//
//

func (lw *LawyerHandler) GetAllLawyers(w http.ResponseWriter, r *http.Request) {
	AllLawyers, err := lw.LawyerStore.GetAllLawyers()
	if err != nil {
		lw.Logger.Printf("Error getting the data from the DB: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"errors": "can't get the data:"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelop{"success": "data retreived:"})
	json.NewEncoder(w).Encode(AllLawyers)
}

func (lw *LawyerHandler) Deletelawyer(w http.ResponseWriter, r *http.Request) {
	IDparams, err := utils.IdReader(r)
	if err != nil {
		lw.Logger.Printf("error getting the id from the param: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"errors": "error getting the id"})
		return
	}
	err = lw.LawyerStore.Deletelawyer(IDparams)
	if err != nil {
		lw.Logger.Printf("Error deleting lawyer: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"errors": "operation failed:"})
		return
	}

	if err == sql.ErrNoRows {
		lw.Logger.Printf("no data with this ID: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"errors": "No data with this ID :"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	utils.WriteJson(w, http.StatusOK, utils.Envelop{"success": "data deleted"})

}

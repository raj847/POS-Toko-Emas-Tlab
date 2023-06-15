package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tim-b/entity"
	"tim-b/service"

	"gorm.io/gorm"
)

type BentukAPI struct {
	bentukService *service.BentukService
}

func NewBentukAPI(
	bentukService *service.BentukService,
) *BentukAPI {
	return &BentukAPI{
		bentukService: bentukService,
	}
}

func (j *BentukAPI) GetAllBentuk(w http.ResponseWriter, r *http.Request) {

	bentuk := r.URL.Query()
	bentukID, foundBentukId := bentuk["bentuk_id"]

	if foundBentukId {
		jID, _ := strconv.Atoi(bentukID[0])
		bentukByID, err := j.bentukService.GetBentukByID(r.Context(), jID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		if bentukByID.ID == 0 {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse("error bentuk not found"))
			return
		}

		WriteJSON(w, http.StatusOK, bentukByID)
		return
	}

	list, err := j.bentukService.GetAllBentuk(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (j *BentukAPI) CreateNewBentuk(w http.ResponseWriter, r *http.Request) {
	var bentuk entity.BentukReq

	err := json.NewDecoder(r.Body).Decode(&bentuk)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid bentuk request"))
		return
	}

	if bentuk.Nama == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid name request"))
		return
	}

	usernameClaims := r.Context().Value("username").(string)

	_, err = j.bentukService.AddBentuk(r.Context(), entity.Bentuk{
		Nama:      bentuk.Nama,
		Kode:      bentuk.Kode,
		CreatedBy: usernameClaims,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"created_by": usernameClaims,
		"message":    "success create new bentuk",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (j *BentukAPI) DeleteBentuk(w http.ResponseWriter, r *http.Request) {
	bentukID := r.URL.Query().Get("bentuk_id")
	jID, _ := strconv.Atoi(bentukID)
	err := j.bentukService.DeleteBentuk(r.Context(), jID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"bentuk_id": jID,
		"message":   "success delete bentuk",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (j *BentukAPI) UpdateBentuk(w http.ResponseWriter, r *http.Request) {
	var bentuk entity.BentukReq

	err := json.NewDecoder(r.Body).Decode(&bentuk)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid bentuk request"))
		return
	}

	usernameClaims := r.Context().Value("username").(string)

	id := r.URL.Query().Get("bentuk_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid bentuk id"))
		return
	}

	bentuks, err := j.bentukService.UpdateBentuk(r.Context(), entity.Bentuk{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		Nama:      bentuk.Nama,
		Kode:      bentuk.Kode,
		UpdatedBy: usernameClaims,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"updated_by": usernameClaims,
		"bentuk_id":  bentuks.ID,
		"message":    "success update bentuk",
	}

	WriteJSON(w, http.StatusOK, response)
}

package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tim-b/entity"
	"tim-b/service"

	"gorm.io/gorm"
)

type JenisAPI struct {
	jenisService *service.JenisService
}

func NewJenisAPI(
	jenisService *service.JenisService,
) *JenisAPI {
	return &JenisAPI{
		jenisService: jenisService,
	}
}

func (j *JenisAPI) GetAllJenis(w http.ResponseWriter, r *http.Request) {

	jenis := r.URL.Query()
	jenisID, foundJenisId := jenis["jenis_id"]

	if foundJenisId {
		jID, _ := strconv.Atoi(jenisID[0])
		jenisByID, err := j.jenisService.GetJenisByID(r.Context(), jID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		if jenisByID.ID == 0 {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse("error jenis not found"))
			return
		}

		WriteJSON(w, http.StatusOK, jenisByID)
		return
	}

	list, err := j.jenisService.GetAllJenis(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (j *JenisAPI) CreateNewJenis(w http.ResponseWriter, r *http.Request) {
	var jenis entity.JenisReq

	err := json.NewDecoder(r.Body).Decode(&jenis)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid jenis request"))
		return
	}

	if jenis.Nama == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid name request"))
		return
	}

	usernameClaims := r.Context().Value("username").(string)

	_, err = j.jenisService.AddJenis(r.Context(), entity.Jenis{
		Nama:      jenis.Nama,
		Kode:      jenis.Kode,
		CreatedBy: usernameClaims,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"created_by": usernameClaims,
		"message":    "success create new jenis",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (j *JenisAPI) DeleteJenis(w http.ResponseWriter, r *http.Request) {
	jenisID := r.URL.Query().Get("jenis_id")
	jID, _ := strconv.Atoi(jenisID)
	err := j.jenisService.DeleteJenis(r.Context(), jID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"jenis_id": jID,
		"message":  "success delete jenis",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (j *JenisAPI) UpdateJenis(w http.ResponseWriter, r *http.Request) {
	var jenis entity.JenisReq

	err := json.NewDecoder(r.Body).Decode(&jenis)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid jenis request"))
		return
	}

	usernameClaims := r.Context().Value("username").(string)

	id := r.URL.Query().Get("jenis_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid jenis id"))
		return
	}

	jeniss, err := j.jenisService.UpdateJenis(r.Context(), entity.Jenis{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		Nama:      jenis.Nama,
		Kode:      jenis.Kode,
		UpdatedBy: usernameClaims,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"updated_by": usernameClaims,
		"jenis_id":   jeniss.ID,
		"message":    "success update jenis",
	}

	WriteJSON(w, http.StatusOK, response)
}

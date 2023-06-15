package api

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"tim-b/entity"
	"tim-b/service"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type InventoryAPI struct {
	inventoryService *service.InventoryService
	minioClient      *minio.Client
}

func NewInventoryAPI(
	inventoryService *service.InventoryService,
	minioClient *minio.Client,
) *InventoryAPI {
	return &InventoryAPI{
		inventoryService: inventoryService,
		minioClient:      minioClient,
	}
}

func (i *InventoryAPI) AddInventory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// get form file from request
	file, header, err := r.FormFile("foto_1")
	if err != nil {
		fmt.Println(file)
		fmt.Println(err.Error())
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid file"))
		return
	}
	defer file.Close()

	file2, header2, err := r.FormFile("foto_2")
	if err != nil {
		fmt.Println(file)
		fmt.Println(err.Error())
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid file"))
		return
	}
	defer file2.Close()

	_, err = i.minioClient.PutObject(r.Context(), "rajendra", header.Filename, file, header.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
		ContentType: "image/jpeg",
	})
	if err != nil {
		log.Println(err)
	}

	_, err = i.minioClient.PutObject(r.Context(), "rajendra", header2.Filename, file2, header2.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
		ContentType: "image/jpeg",
	})
	if err != nil {
		log.Println(err)
	}

	fileName1 := fmt.Sprintf("https://is3.cloudhost.id/rajendra/%s", header.Filename)
	fileName2 := fmt.Sprintf("https://is3.cloudhost.id/rajendra/%s", header2.Filename)

	// get form value from request

	idjenis := r.FormValue("id_jenis_barang")
	idbentuk := r.FormValue("id_jenis_barang")
	namabarang := r.FormValue("nama_barang")
	berat := r.FormValue("berat")
	kadar := r.FormValue("kadar")
	hargajual := r.FormValue("harga_jual")
	catatan := r.FormValue("catatan")

	rand.Seed(time.Now().UnixNano())

	// Menghasilkan angka acak 5 digit.
	randomNum := rand.Intn(90000) + 10000

	// Format angka acak menjadi string dengan 5 digit.
	randomNumStr := fmt.Sprintf("%05d", randomNum)

	idj, _ := strconv.Atoi(idjenis)
	idb, _ := strconv.Atoi(idbentuk)
	brt, _ := strconv.ParseFloat(berat, 64)
	kdr, _ := strconv.ParseFloat(kadar, 64)
	hj, _ := strconv.ParseFloat(hargajual, 64)

	usernameClaims := r.Context().Value("username").(string)

	inventory := entity.InventoryReq{
		IDJenis:    idj,
		IDBentuk:   idb,
		NamaBarang: namabarang,
		Berat:      brt,
		Kadar:      kdr,
		HargaJual:  hj,
		Notes:      catatan,
		PhotoURL1:  fileName1,
		PhotoURL2:  fileName2,
	}

	err = r.ParseMultipartForm(4096) // parsing request dengan size maksimal 4096 bytes
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = i.inventoryService.AddInventory(r.Context(), entity.Inventory{
		IDJenis:    inventory.IDJenis,
		IDBentuk:   inventory.IDBentuk,
		NamaBarang: inventory.NamaBarang,
		KodeBarang: randomNumStr,
		Berat:      inventory.Berat,
		Kadar:      inventory.Kadar,
		HargaJual:  inventory.HargaJual,
		Notes:      inventory.Notes,
		PhotoURL1:  inventory.PhotoURL1,
		PhotoURL2:  inventory.PhotoURL2,
		CreatedBy:  usernameClaims,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"message": "add inventory success",
	}

	WriteJSON(w, http.StatusCreated, response)
}
func (i *InventoryAPI) GetAllInventory(w http.ResponseWriter, r *http.Request) {

	inv := r.URL.Query()
	invID, foundJenisId := inv["inv_id"]

	if foundJenisId {
		jID, _ := strconv.Atoi(invID[0])
		invByID, err := i.inventoryService.ReadInventoryID(r.Context(), jID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		if invByID.ID == 0 {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse("error jenis not found"))
			return
		}

		WriteJSON(w, http.StatusOK, invByID)
		return
	}

	list, err := i.inventoryService.ReadInventory()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, list)
}

func (i *InventoryAPI) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	invID := r.URL.Query().Get("inv_id")
	jID, _ := strconv.Atoi(invID)
	err := i.inventoryService.DeleteInventory(r.Context(), jID)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"inv_id":  jID,
		"message": "success delete inventory",
	}

	WriteJSON(w, http.StatusOK, response)
}

func (i *InventoryAPI) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// get form file from request
	file, header, err := r.FormFile("foto_1")
	if err != nil {
		fmt.Println(file)
		fmt.Println(err.Error())
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid file"))
		return
	}
	defer file.Close()

	file2, header2, err := r.FormFile("foto_2")
	if err != nil {
		fmt.Println(file)
		fmt.Println(err.Error())
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid file"))
		return
	}
	defer file2.Close()

	_, err = i.minioClient.PutObject(r.Context(), "rajendra", header.Filename, file, header.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
		ContentType: "image/jpeg",
	})
	if err != nil {
		log.Println(err)
	}

	_, err = i.minioClient.PutObject(r.Context(), "rajendra", header2.Filename, file2, header2.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
		ContentType: "image/jpeg",
	})
	if err != nil {
		log.Println(err)
	}

	fileName1 := fmt.Sprintf("https://is3.cloudhost.id/rajendra/%s", header.Filename)
	fileName2 := fmt.Sprintf("https://is3.cloudhost.id/rajendra/%s", header2.Filename)

	// get form value from request

	idjenis := r.FormValue("id_jenis_barang")
	idbentuk := r.FormValue("id_jenis_barang")
	namabarang := r.FormValue("nama_barang")
	berat := r.FormValue("berat")
	kadar := r.FormValue("kadar")
	hargajual := r.FormValue("harga_jual")
	catatan := r.FormValue("catatan")

	rand.Seed(time.Now().UnixNano())

	// Menghasilkan angka acak 5 digit.
	randomNum := rand.Intn(90000) + 10000

	// Format angka acak menjadi string dengan 5 digit.
	randomNumStr := fmt.Sprintf("%05d", randomNum)

	idj, _ := strconv.Atoi(idjenis)
	idb, _ := strconv.Atoi(idbentuk)
	brt, _ := strconv.ParseFloat(berat, 64)
	kdr, _ := strconv.ParseFloat(kadar, 64)
	hj, _ := strconv.ParseFloat(hargajual, 64)

	usernameClaims := r.Context().Value("username").(string)
	id := r.URL.Query().Get("inv_id")
	idInt, err := strconv.Atoi(id)

	inventory := entity.InventoryReq{
		IDJenis:    idj,
		IDBentuk:   idb,
		NamaBarang: namabarang,
		Berat:      brt,
		Kadar:      kdr,
		HargaJual:  hj,
		Notes:      catatan,
		PhotoURL1:  fileName1,
		PhotoURL2:  fileName2,
	}

	err = r.ParseMultipartForm(4096) // parsing request dengan size maksimal 4096 bytes
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = i.inventoryService.UpdateInventory(r.Context(), entity.Inventory{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		IDJenis:    inventory.IDJenis,
		IDBentuk:   inventory.IDBentuk,
		NamaBarang: inventory.NamaBarang,
		KodeBarang: randomNumStr,
		Berat:      inventory.Berat,
		Kadar:      inventory.Kadar,
		HargaJual:  inventory.HargaJual,
		Notes:      inventory.Notes,
		PhotoURL1:  inventory.PhotoURL1,
		PhotoURL2:  inventory.PhotoURL2,
		UpdatedBy:  usernameClaims,
	})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"id":      idInt,
		"message": "update inventory success",
	}

	WriteJSON(w, http.StatusCreated, response)
}

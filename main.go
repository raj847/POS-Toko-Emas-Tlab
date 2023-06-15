package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tim-b/handler/api"
	"tim-b/middleware"
	"tim-b/repository"
	"tim-b/service"
	"tim-b/utils"

	"github.com/rs/cors"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler      *api.UserAPI
	JenisAPIHandler     *api.JenisAPI
	BentukAPIHandler    *api.BentukAPI
	InventoryAPIHandler *api.InventoryAPI
}

func main() {
	err := os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/tim-b")
	if err != nil {
		log.Fatalf("cannot set env: %v", err)
	}

	mux := http.NewServeMux()

	err = utils.ConnectDB()
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	db := utils.GetDBConnection()
	mux = RunServer(db, mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}

func RunServer(db *gorm.DB, mux *http.ServeMux) *http.ServeMux {
	minioClientConn, err := service.NewMinioClient()
	if err != nil {
		log.Fatalf("cannot connect to minio: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	jenisRepo := repository.NewJenisRepository(db)
	bentukRepo := repository.NewBentukRepository(db)
	invRepo := repository.NewInventoryRepository(db)

	userService := service.NewUserService(userRepo)
	jenisService := service.NewJenisService(jenisRepo)
	bentukService := service.NewBentukService(bentukRepo)
	invService := service.NewInventoryService(invRepo)

	userAPIHandler := api.NewUserAPI(userService)
	jenisAPIHandler := api.NewJenisAPI(jenisService)
	bentukAPIHandler := api.NewBentukAPI(bentukService)
	invAPIHandler := api.NewInventoryAPI(invService, minioClientConn)

	apiHandler := APIHandler{
		UserAPIHandler:      userAPIHandler,
		JenisAPIHandler:     jenisAPIHandler,
		BentukAPIHandler:    bentukAPIHandler,
		InventoryAPIHandler: invAPIHandler,
	}

	//USER
	MuxRoute(mux, "POST", "/tim-b/ver1/user/add-anggota", middleware.Post(
		middleware.Auth(
			http.HandlerFunc(
				apiHandler.UserAPIHandler.AddAnggota))))
	MuxRoute(mux, "POST", "/tim-b/ver1/user/register", middleware.Post(
		http.HandlerFunc(
			apiHandler.UserAPIHandler.Register)))
	MuxRoute(mux, "POST", "/tim-b/ver1/user/login", middleware.Post(
		http.HandlerFunc(
			apiHandler.UserAPIHandler.UserLogin)))
	MuxRoute(mux, "GET", "/tim-b/ver1/user/get-anggota", middleware.Get(
		middleware.Auth(
			http.HandlerFunc(
				apiHandler.UserAPIHandler.GetAllUsers))),
		"?user_id=", "&created_by=", "&nama=",
	)
	MuxRoute(mux, "PUT", "/tim-b/ver1/user/update-anggota", middleware.Put(
		middleware.Auth(
			http.HandlerFunc(
				apiHandler.UserAPIHandler.UpdateAnggota))),
		"?user_id=",
	)
	//JENIS
	MuxRoute(mux, "POST", "/tim-b/ver1/jenis/create",
		middleware.Post(
			middleware.Auth(
				http.HandlerFunc(apiHandler.JenisAPIHandler.CreateNewJenis))))

	MuxRoute(mux, "GET", "/tim-b/ver1/jenis",
		middleware.Get(
			middleware.Auth(
				http.HandlerFunc(apiHandler.JenisAPIHandler.GetAllJenis),
			),
		),
		"?jenis_id=",
	)

	MuxRoute(mux, "PUT", "/tim-b/ver1/jenis/update",
		middleware.Put(
			middleware.Auth(
				http.HandlerFunc(apiHandler.JenisAPIHandler.UpdateJenis))),
		"?jenis_id=",
	)

	MuxRoute(mux, "DELETE", "/tim-b/ver1/jenis/delete",
		middleware.Delete(
			middleware.Auth(
				http.HandlerFunc(apiHandler.JenisAPIHandler.DeleteJenis))),
		"?jenis_id=",
	)
	//BENTUK
	MuxRoute(mux, "POST", "/tim-b/ver1/bentuk/create",
		middleware.Post(
			middleware.Auth(
				http.HandlerFunc(apiHandler.BentukAPIHandler.CreateNewBentuk))))

	MuxRoute(mux, "GET", "/tim-b/ver1/bentuk",
		middleware.Get(
			middleware.Auth(
				http.HandlerFunc(apiHandler.BentukAPIHandler.GetAllBentuk),
			),
		),
		"?bentuk_id=",
	)

	MuxRoute(mux, "PUT", "/tim-b/ver1/bentuk/update",
		middleware.Put(
			middleware.Auth(
				http.HandlerFunc(apiHandler.BentukAPIHandler.UpdateBentuk))),
		"?bentuk_id=",
	)

	MuxRoute(mux, "DELETE", "/tim-b/ver1/bentuk/delete",
		middleware.Delete(
			middleware.Auth(
				http.HandlerFunc(apiHandler.BentukAPIHandler.DeleteBentuk))),
		"?bentuk_id=",
	)
	//INVENTORY
	MuxRoute(mux, "POST", "/tim-b/ver1/inv/create",
		middleware.Post(
			middleware.Auth(
				http.HandlerFunc(apiHandler.InventoryAPIHandler.AddInventory))))

	MuxRoute(mux, "GET", "/tim-b/ver1/inv",
		middleware.Get(
			middleware.Auth(
				http.HandlerFunc(apiHandler.InventoryAPIHandler.GetAllInventory),
			),
		),
		"?inv_id=",
	)

	MuxRoute(mux, "PUT", "/tim-b/ver1/inv/update",
		middleware.Put(
			middleware.Auth(
				http.HandlerFunc(apiHandler.InventoryAPIHandler.UpdateInventory))),
		"?inv_id=",
	)

	MuxRoute(mux, "DELETE", "/tim-b/ver1/inv/delete",
		middleware.Delete(
			middleware.Auth(
				http.HandlerFunc(apiHandler.InventoryAPIHandler.DeleteInventory))),
		"?inv_id=",
	)

	return mux

}

func MuxRoute(mux *http.ServeMux, method string, path string, handler http.Handler, opt ...string) {
	if len(opt) > 0 {
		fmt.Printf("[%s]: %s %v \n", method, path, opt)
	} else {
		fmt.Printf("[%s]: %s \n", method, path)
	}

	mux.Handle(path, handler)
}

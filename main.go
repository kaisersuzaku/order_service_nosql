package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kaisersuzaku/order_service_nosql/handlers"

	"github.com/kaisersuzaku/order_service_nosql/services"

	"github.com/kaisersuzaku/order_service_nosql/repo"

	"github.com/kaisersuzaku/order_service_nosql/utils"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	utils.GetConfig(dir + "/conf.json")
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	opr := repo.BuildProductRepo(utils.BuildDataStore())
	ops := services.BuildOrderProductService(opr)
	oph := handlers.BuildOrderProductHandler(ops)

	r.Post("/order-product", oph.OrderProduct)

	http.ListenAndServe(":7789", r)
}

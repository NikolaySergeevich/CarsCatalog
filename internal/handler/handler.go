package handler

import (
	// "fmt"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"testcar/internal/database"

	// "testcar/internal/database/car"
	"testcar/internal/env"
	"testcar/pkg/api/carapi"
)

var _ carapi.ServerInterface = (*Handler)(nil)
/*
func NewHandler(db car.Repository) *Handler{
	return &Handler{db: db}
}

type Handler struct {
	db car.Repository
}
*/
func NewHandler(ctx context.Context, env env.Env) *Handler{
	return &Handler{ctx: ctx, env: &env}
}

type Handler struct{
	ctx context.Context
	env *env.Env
}
func (h *Handler) GetCars(w http.ResponseWriter, r *http.Request, params carapi.GetCarsParams) {
	//to do
	fmt.Println("куку")
	w.WriteHeader(http.StatusOK)
}
func (h *Handler) PostCars(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sdfsdfsdf")
	serverPort := 8789 //порт сервера для внешнего api, у которого будет запрашиваться инф. о машине.
	//to do
	//итерируемся по списку, который передали в json при запросе. И на каждой итерации выполняем запрос ко внешнему api
	//собираем список из CreateCar и передаём в БД. Но лучше тут использовать очереди
	listRegNums := make(map[string][]string)
	listCreateCar := []database.CreateCar{}

	content, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil{
		h.env.Logger.Error("Read Body from other service", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.Unmarshal(content, &listRegNums); err != nil {
		h.env.Logger.Error("json.Unmarshal from other service", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	crCar := database.CreateCar{}
	for _, regNum := range listRegNums {
		requestURL := fmt.Sprintf("http://localhost:%d/api/v1/info/%s", serverPort, regNum)
		res, err := http.Get(requestURL)
		if err != nil {
			h.env.Logger.Error("http get other service", slog.Any("err", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewDecoder(res.Body).Decode(&crCar); err != nil {
			h.env.Logger.Error("Decoder Car from other service", slog.Any("err", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		listCreateCar = append(listCreateCar, crCar)
	}
	if err := h.env.AutoRepository.Create(h.ctx, listCreateCar); err != nil {
		h.env.Logger.Error("Add cars db", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteCarsId(w http.ResponseWriter, r *http.Request, id string) {

}

func (h *Handler) PutCarsId(w http.ResponseWriter, r *http.Request, id string) {

}


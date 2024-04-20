package handler

import (
	// "fmt"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strconv"
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
func NewHandler(ctx context.Context, env env.Env) *Handler {
	return &Handler{ctx: ctx, env: &env}
}

type Handler struct {
	ctx context.Context
	env *env.Env
}

func (h *Handler) GetCars(w http.ResponseWriter, r *http.Request, params carapi.GetCarsParams) {
	//to do
	paramss := r.URL.Query()

	// Обработка параметров limit и page
	limit := 20 // По умолчанию
	if limitStr := paramss.Get("limit"); limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	page := 1 // По умолчанию
	if pageStr := paramss.Get("page"); pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	query := `
		SELECT id, mark, model, color, yearCr, regNums, ownerCar, created_at, updated_at
		FROM cars 
		WHERE TRUE`
	args := []interface{}{}

	mark := paramss.Get("mark")
	if mark != "" {
		query += " AND mark = $1"
		args = append(args, mark)
	} else {
		// args = append(args, "qw")
	}
	model := paramss.Get("model")
	if model != "" {
		query += " AND model = $2"
		args = append(args, model)
	} else {
		// args = append(args, "qw")
	}
	regNums := paramss.Get("regNums");
	if regNums != "" {
		query += " AND regNums = $3"
		args = append(args, regNums)
	} else{
		// args = append(args, "qw")
	}
	owner := paramss.Get("ownerCar");
	if owner != "" {
		query += " AND ownerCar = $4"
		args = append(args, owner)
	} else{
		// args = append(args, "qw")
	}
	color := paramss.Get("color");
	if color != "" {
		query += " AND color = $2"
		args = append(args, color)
	} else{
		// args = append(args, "qw")
	}
	year := paramss.Get("yearCr");
	if year != "" {
		query += " AND yearCr::text = $6"
		args = append(args, year)
	} else{
		// args = append(args, "qw")
	}

	cars, err := h.env.AutoRepository.TakeCars(h.ctx, query, args)
	if err != nil {
		h.env.Logger.Error("DB query:", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	startIndex := (page - 1) * limit
	endIndex := startIndex + limit
	if endIndex > len(cars) {
		endIndex = len(cars)
	}
	paginatedCars := cars[startIndex:endIndex]

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paginatedCars); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *Handler) PostCars(w http.ResponseWriter, r *http.Request) {
	serverPort := 8789 //порт сервера для внешнего api, у которого будет запрашиваться инф. о машине.
	//to do
	//итерируемся по списку, который передали в json при запросе. И на каждой итерации выполняем запрос ко внешнему api
	//собираем список из CreateCar и передаём в БД. Но лучше тут использовать очереди
	listRegNums := make(map[string][]string)
	listCreateCar := []database.CreateCar{}

	content, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
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

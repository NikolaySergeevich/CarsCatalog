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

	"github.com/google/uuid"
)

var _ carapi.ServerInterface = (*Handler)(nil)

func NewHandler(ctx context.Context, env env.Env) *Handler {
	return &Handler{ctx: ctx, env: &env}
}

type Handler struct {
	ctx context.Context
	env *env.Env
}

func (h *Handler) GetCars(w http.ResponseWriter, r *http.Request, params carapi.GetCarsParams) {
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
	num := 1

	if mark := paramss.Get("mark"); mark != "" {
		query += fmt.Sprintf(" AND mark = $%s", strconv.Itoa(num))
		num = num + 1
		args = append(args, mark)
	}
	if model := paramss.Get("model"); model != "" {
		query += fmt.Sprintf(" AND model = $%s", strconv.Itoa(num))
		num = num + 1
		args = append(args, model)
	}
	if regNums := paramss.Get("regNums"); regNums != "" {
		query += fmt.Sprintf(" AND regNums = $%s", strconv.Itoa(num))
		num = num + 1
		args = append(args, regNums)
	}
	if owner := paramss.Get("ownerCar"); owner != "" {
		query += fmt.Sprintf(" AND ownerCar = $%s", strconv.Itoa(num))
		num = num + 1
		args = append(args, owner)
	}
	if color := paramss.Get("color"); color != "" {
		query += fmt.Sprintf(" AND color = $%s", strconv.Itoa(num))
		num = num + 1
		args = append(args, color)
	}
	if year := paramss.Get("yearCr"); year != "" {
		query += fmt.Sprintf(" AND yearCr::text = $%s", strconv.Itoa(num))
		num = num + 1
		args = append(args, year)
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
	RegNums := make(map[string][]string)
	CreateCar := []database.CreateCar{}

	content, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		h.env.Logger.Error("Read Body from other service", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.Unmarshal(content, &RegNums); err != nil {
		h.env.Logger.Error("json.Unmarshal from other service", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(RegNums)
	if _, ok := RegNums["regNums"]; !ok{
		h.env.Logger.Error("Invalid request specified json", slog.Any("err", "You need 'regNums': ['X123XX150']"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	crCar := database.CreateCar{}
	for _, regNum := range RegNums["regNums"] {
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
		CreateCar = append(CreateCar, crCar)
	}
	if err := h.env.AutoRepository.Create(h.ctx, CreateCar); err != nil {
		h.env.Logger.Error("Add cars db", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteCarsId(w http.ResponseWriter, r *http.Request, id string) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		h.env.Logger.Error("Pars id", slog.Any("err", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.env.AutoRepository.DeleteCarsId(h.ctx, parsedUUID); err != nil {
		h.env.Logger.Error("Delete car", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PutCarsId(w http.ResponseWriter, r *http.Request, id string) {

	content, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		h.env.Logger.Error("Read Body from other service", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	newCar := carapi.PutCarsIdJSONBody{}
	if err := json.Unmarshal(content, &newCar); err != nil {
		h.env.Logger.Error("json.Unmarshal from other service", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	args := []interface{}{}
	query := `UPDATE cars SET`
	num := 1
	if mark := *newCar.Mark; mark != "" {
		query += fmt.Sprintf(" mark = $%s,", strconv.Itoa(num))
		num = num + 1
		args = append(args, mark)
	}
	if model := *newCar.Model; model != "" {
		query += fmt.Sprintf(" model = $%s,", strconv.Itoa(num))
		num = num + 1
		args = append(args, model)
	}
	if color := *newCar.Color; color != "" {
		query += fmt.Sprintf(" color = $%s,", strconv.Itoa(num))
		num = num + 1
		args = append(args, color)
	}
	if owner := *newCar.Owner; owner != "" {
		query += fmt.Sprintf(" ownerCar = $%s,", strconv.Itoa(num))
		num = num + 1
		args = append(args, owner)
	}
	//убирается последняя запятая
	query = query[:len(query)-1] + fmt.Sprintf(" WHERE id = $%s", strconv.Itoa(num))
	args = append(args, id)

	fmt.Println(query)
	fmt.Println(args...)
	if err := h.env.AutoRepository.UpdateCar(h.ctx, query, args); err != nil {
		h.env.Logger.Error("DB query:", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

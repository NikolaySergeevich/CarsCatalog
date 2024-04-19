package handler

import (
	// "fmt"
	"fmt"
	"net/http"
	"testcar/internal/database/car"
	"testcar/pkg/api/carapi"
)

var _ carapi.ServerInterface = (*Handler)(nil)

func NewHandler(db car.Repository) *Handler{
	return &Handler{db: db}
}

type Handler struct {
	db car.Repository
}
func (h *Handler) GetCars(w http.ResponseWriter, r *http.Request, params carapi.GetCarsParams) {
	//to do
	fmt.Println("куку")
	w.WriteHeader(http.StatusOK)
}
func (h *Handler) PostCars(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("sdfsdfsdf")
	// serverPort := 8789 //порт сервера для внешнего api, у которого будет запрашиваться инф. о машине.
	// //to do
	// //итерируемся по списку, который передали в json при запросе. И на каждой итерации выполняем запрос ко внешнему api
	// //собираем список из CreateCar и передаём в БД. Но лучше тут использовать очереди
	// listRegNums := []string{}
	// listCreateCar := []database.CreateCar{}
	// if err := json.NewDecoder(r.Body).Decode(&listRegNums); err != nil {
	// 	h.setup.Logger.Error("json.NewDecoder Decode from other service", slog.Any("err", err))
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// fmt.Println(listRegNums)

	// crCar := database.CreateCar{}
	// for _, regNum := range listRegNums {
	// 	requestURL := fmt.Sprintf("http://localhost:%d/api/v1/info/%s", serverPort, regNum)
	// 	res, err := http.Get(requestURL)
	// 	if err != nil {
	// 		h.setup.Logger.Error("http get other service", slog.Any("err", err))
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	if err := json.NewDecoder(res.Body).Decode(&crCar); err != nil {
	// 		h.setup.Logger.Error("Decoder Car from other service", slog.Any("err", err))
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	listCreateCar = append(listCreateCar, crCar)
	// }
	// _, err := h.setup.AutoRepository.Create(h.ctx, listCreateCar)
	// if err != nil {
	// 	h.setup.Logger.Error("Add cars db", slog.Any("err", err))
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteCarsId(w http.ResponseWriter, r *http.Request, id string) {

}

func (h *Handler) PutCarsId(w http.ResponseWriter, r *http.Request, id string) {

}


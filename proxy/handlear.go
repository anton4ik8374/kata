package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	geoApi *GeoService
}

func NewHandler() *Handler {
	var gs = NewGeoService("da8e86a1acc10757fc82b72fc57ffc8644191336", "b79dc044fc402dff6276c63fc1d8a6aad1b5bc1e")
	return &Handler{
		geoApi: gs,
	}
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	var sarc SearchRequest
	err := json.NewDecoder(r.Body).Decode(&sarc) // считываем приходящий json из *http.Request в структуру SearchRequest
	if err != nil {                              // в случае ошибки отправляем ошибку Bad request code 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return // не забываем прекратить обработку нашего handler (ручки)
	}

	var result SearchResponse
	//Сервис должен возвращать статус 500 в случае, если сервис https://dadata.ru не доступен.
	if h.geoApi == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return //
	}

	resReq, err := h.geoApi.AddressSearch(sarc.Query)
	if err != nil { // в случае ошибки отправляем ошибку Bad request code 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return // не забываем прекратить обработку нашего handler (ручки)
	}
	result.Addresses = resReq

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err = json.NewEncoder(w).Encode(result) // записываем результат GeocodeResponse json в http.SearchResponse

	if err != nil { // отправляем 500 ошибку в случае неудачи
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // не забываем прекратить обработку нашего handler (ручки)
	}

}

func (h *Handler) geocode(w http.ResponseWriter, r *http.Request) {
	var geo GeocodeRequest
	err := json.NewDecoder(r.Body).Decode(&geo) // считываем приходящий json из *http.Request в структуру GeocodeRequest
	if err != nil {                             // в случае ошибки отправляем ошибку Bad request code 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return // не забываем прекратить обработку нашего handler (ручки)
	}

	var result GeocodeResponse
	//Сервис должен возвращать статус 500 в случае, если сервис https://dadata.ru не доступен.
	if h.geoApi == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return //
	}

	resReq, err := h.geoApi.GeoCode(geo.Lat, geo.Lng)
	if err != nil { // в случае ошибки отправляем ошибку Bad request code 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return // не забываем прекратить обработку нашего handler (ручки)
	}
	result.Addresses = resReq

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err = json.NewEncoder(w).Encode(result) // записываем результат GeocodeResponse json в http.ResponseWriter

	if err != nil { // отправляем 500 ошибку в случае неудачи
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // не забываем прекратить обработку нашего handler (ручки)
	}
}

func (h *Handler) swaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var htmlB, err = os.ReadFile("./swagger/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl, err := template.New("swagger").Parse(string(htmlB))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, struct {
		Time int64
	}{
		Time: time.Now().Unix(),
	})
	if err != nil {
		return
	}
}

func (h *Handler) swaggerGET(w http.ResponseWriter, r *http.Request) {
	var htmlB, err = os.ReadFile("./swagger/swagger.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(htmlB)
}

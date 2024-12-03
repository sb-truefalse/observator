package main

import (
	"net/http"
	"strconv"
)

type MemStorageFuncs interface {
	setCounter()
	setGauge()
}

type MemStorage struct {
	counter map[string]int64
	gauge   map[string]float64
}

func (memStorage MemStorage) setCounter(name string, value int64) {
	memStorage.counter[name] += value
}

func (memStorage MemStorage) setGauge(name string, value float64) {
	memStorage.gauge[name] = value
}

func updateMetric(res http.ResponseWriter, req *http.Request) {
	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metricValue := req.PathValue("metricValue")

	// Check method
	if req.Method == http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check type
	if metricType != "counter" && metricType != "gauge" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check name
	if metricName == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if metricType == "counter" {
		// Check and set value
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err == nil {
			currentMemStorage.setCounter(metricName, value)
		} else {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		// Check and set value
		value, err := strconv.ParseFloat(metricValue, 64)
		if err == nil {
			currentMemStorage.setGauge(metricName, value)
		} else {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	res.Header().Set("content-type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
}

var currentMemStorage MemStorage

func main() {
	// Init MemStorage
	currentMemStorage.counter = make(map[string]int64)
	currentMemStorage.gauge = make(map[string]float64)

	// Init server
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/{metricType}/{metricName}/{metricValue}`, updateMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

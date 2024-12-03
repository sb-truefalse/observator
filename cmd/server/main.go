package main

import (
	"net/http"
	"strconv"
)

type MemStorageFuncs interface {
	getCounter()
	getGauge()
	setCounter()
	setGauge()
}

type MemStorage struct {
	counter map[string]int64
	gauge   map[string]float64
}

func (memStorage MemStorage) getCounter(name string) int64 {
	return memStorage.counter[name]
}

func (memStorage MemStorage) getGauge(name string) float64 {
	return memStorage.gauge[name]
}

func (memStorage MemStorage) setCounter(name string, value int64) {
	memStorage.counter[name] += value
}

func (memStorage MemStorage) setGauge(name string, value float64) {
	memStorage.gauge[name] = value
}

func updateMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		metricType := req.PathValue("metricType")
		metricName := req.PathValue("metricName")
		metricValue := req.PathValue("metricValue")

		if metricType == "counter" || metricType == "gauge" {
			if metricName == "" {
				res.WriteHeader(http.StatusNotFound)
				return
			} else {

				if metricType == "counter" {
					_, err := strconv.Atoi(metricValue)
					if err == nil {
						//
					} else {
						res.WriteHeader(http.StatusBadRequest)
						return
					}
				} else {
					_, err := strconv.ParseFloat(metricValue, 64)
					if err == nil {
						//
					} else {
						res.WriteHeader(http.StatusBadRequest)
						return
					}
				}

				res.Header().Set("content-type", "text/plain; charset=utf-8")
				res.WriteHeader(http.StatusOK)
			}
		} else {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		return
	} else {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/{metricType}/{metricName}/{metricValue}`, updateMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

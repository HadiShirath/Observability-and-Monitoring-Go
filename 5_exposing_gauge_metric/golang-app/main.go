package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "golang_app",
		Name:      "http_request_count",
	}, []string{"method", "path", "code"},
	)

	metricGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "golang_app",
		Name:      "active_users",
	}, []string{"country_id", "city_id"})
)

func init() {
	prometheus.MustRegister(metricCounter)
	prometheus.MustRegister(metricGauge)
}

func main() {
	println("starting http server ...")

	ctx, cancelSimulationJob := context.WithCancel(context.Background())
	defer cancelSimulationJob()

	sigCH := make(chan os.Signal, 1)
	signal.Notify(sigCH, syscall.SIGINT, syscall.SIGTERM)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				println("stopping job simulation")
				return
			default:
				log.Printf("sending metrics at %s\n", time.Now().Format(time.RFC3339Nano))
				metricGauge.WithLabelValues("ID", "JAK").Add(float64(rand.Intn(1000)))
				time.Sleep(300 * time.Millisecond)
			}
		}
	}(ctx)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`pong`))
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(450 * time.Millisecond)
		responseData := ResponseData{
			Data: Data{
				OrderID: rand.Intn(1000),
			},
		}

		// to simulate successfull or failing http status code
		if responseData.Data.OrderID%2 == 0 {
			metricCounter.WithLabelValues("POST", "/orders", "200")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responseData)
		} else {
			metricCounter.WithLabelValues("POST", "/orders", "500")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			responseData.Error = Error{
				Message:     "failure, this is simulated by test devs",
				Description: "error",
			}
			json.NewEncoder(w).Encode(responseData)
		}

	})

	http.ListenAndServe(":1000", nil)
	println("stopped http server ...")

}

type (
	ResponseData struct {
		Data  Data  `json:"data"`
		Error Error `json:"error,omitempty"`
	}

	Data struct {
		OrderID int `json:"order_id"`
	}

	Error struct {
		Message     string `json:"message,omitempty"`
		Description string `json:"description,omitempty"`
	}
)

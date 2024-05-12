package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HPMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hp",
		Help: "The total number of Health Points",
	})
	ArmorMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "armor",
		Help: "Persent of of Armor",
	})
)

type Jedi struct {
	Name         string
	HP           int // Жизни
	Force        int
	ClonePlatoon string
}

func InitJedi() *Jedi {
	clone := Jedi{
		Name:         os.Getenv("HOSTNAME"),
		HP:           GetEnv("HP", 100),
		Force:        GetEnv("FORCE", 10),
		ClonePlatoon: GetUrl("CLONE_PLATOON", "http://clone-platoon"),
	}
	return &clone
}

func (j *Jedi) recordMetrics() {
	go func() {
		for {
			HPMetric.Set(float64(j.HP))
			time.Sleep(2 * time.Second)
		}
	}()
}

func (j *Jedi) attack() {
	resp, err := http.Get(fmt.Sprintf("%s/damage", j.ClonePlatoon))
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			log.Warn().Msg(err.Error())
			return
		}
		log.Error().Msg(err.Error())
		return
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	log.Info().Msg(string(bodyBytes))
}

func (j *Jedi) AttackClones() {
	go func() {
		for {
			j.attack()
			time.Sleep(2 * time.Second)
		}
	}()
}

func (j *Jedi) HtmlTemplate(w http.ResponseWriter, req *http.Request) {
	tmplt := template.Must(template.ParseFiles("jedi/static/index.html"))
	tmplt.Execute(w, j)
}

func Icon(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := os.ReadFile("jedi/static/img/icon.png")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func main() {

	jedi := InitJedi()
	log.Info().Msgf("Hello, my name is '%s'", jedi.Name)

	jedi.AttackClones()

	r := prometheus.NewRegistry()
	r.MustRegister(HPMetric, ArmorMetric)
	jedi.recordMetrics()

	http.HandleFunc("/", jedi.HtmlTemplate)
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if jedi.HP == 0 {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("error: jedi died")))
		} else {
			w.WriteHeader(200)
			w.Header().Add("name", jedi.Name)
			w.Write([]byte("ok"))
		}
	})
	http.HandleFunc("/icon", Icon)

	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

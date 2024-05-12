package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
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

type Clone struct {
	Name     string
	HP       int     // Жизни
	Accuracy int     // Точность
	Armor    float64 // Броня
}

func InitClone() *Clone {
	clone := Clone{
		Name:     os.Getenv("HOSTNAME"),
		HP:       GetEnv("HP", 100),
		Accuracy: GetEnv("ACCURANCY", 5),
		Armor:    GetEnv("ARMOR", 1.),
	}
	return &clone
}

func (c *Clone) recordMetrics() {
	go func() {
		for {
			HPMetric.Set(float64(c.HP))
			ArmorMetric.Set(c.Armor)
			time.Sleep(2 * time.Second)
		}
	}()
}

func (c *Clone) Damage(w http.ResponseWriter, req *http.Request) {
	max := 100
	armor := 1.0 - c.Armor
	damage := float64(rand.Intn(max))
	if armor > 0 {
		c.HP -= int(damage * armor)
		if c.HP < 0 {
			c.HP = 0
			log.Fatal().Msg("clone died")
		}
	}
	c.Armor -= damage * 0.01
	if c.Armor < 0 {
		c.Armor = 0
	}
	log.Info().
		Str("armor", fmt.Sprintf("%d%%", int(c.Armor*100))).
		Int("hp", c.HP).
		Msgf("get damage from %s", req.RemoteAddr)
	fmt.Fprintf(w, "%s: %d", c.Name, c.HP)
}

func (c *Clone) HtmlTemplate(w http.ResponseWriter, req *http.Request) {
	tmplt := template.Must(template.ParseFiles("clone/static/index.html"))
	tmplt.Execute(w, c)
}

func Icon(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := os.ReadFile("clone/static/img/icon.png")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func main() {

	clone := InitClone()
	log.Info().Msgf("Hello, my name is '%s'", clone.Name)

	r := prometheus.NewRegistry()
	r.MustRegister(HPMetric, ArmorMetric)
	clone.recordMetrics()

	http.HandleFunc("/", clone.HtmlTemplate)
	http.HandleFunc("/damage", clone.Damage)
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if clone.HP == 0 {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("error: clone died")))
		} else {
			w.WriteHeader(200)
			w.Write([]byte("ok"))
		}
	})
	http.HandleFunc("/icon", Icon)

	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}

}

package server

import (
	"bytes"
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
)

//go:embed ui/ctx-dashboard/dist/ctx-dashboard/assets/*
//go:embed ui/ctx-dashboard/dist/ctx-dashboard/*.html
//go:embed ui/ctx-dashboard/dist/ctx-dashboard/*.ico
var staticFiles embed.FS

func spaHandler(content fs.FS, fsHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URL:", r.URL.Path)
		_, err := fs.Stat(content, r.URL.Path[1:])
		if err != nil {
			f, err := content.Open("index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			data, err := io.ReadAll(f)
			if err != nil {
				http.Error(w, "error reading index.html", http.StatusInternalServerError)
				return
			}

			http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(data))
			return
		}

		fsHandler.ServeHTTP(w, r)
	}
}

func contextList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	mgr := ctx.CreateManager()

	json.NewEncoder(w).Encode(mgr.ListJson2())
}

func currentContext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	mgr := ctx.CreateManager()

	mgr.ContextStore.Read(func(s *ctx_model.State) error {
		if s.CurrentId != "" {
			json.NewEncoder(w).Encode(s.Contexts[s.CurrentId])
		} else {
			json.NewEncoder(w).Encode(nil)
		}
		return nil
	})

}

func createAndSwitchContext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	mgr := ctx.CreateManager()

	var p createAndSwitchRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	mgr.CreateIfNotExistsAndSwitch(util.GenerateId(p.Description), p.Description)
}

func updateInterval(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	mgr := ctx.CreateManager()

	var p EditIntervalRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	mgr.EditContextInterval(p.Id, p.IntervalId, p.Start, p.End)
}

func roundDuration(d time.Duration, unit string) time.Duration {
	switch unit {
	case "nanosecond":
		return d.Round(time.Nanosecond)
	case "microsecond":
		return d.Round(time.Microsecond)
	case "millisecond":
		return d.Round(time.Millisecond)
	case "second":
		return d.Round(time.Second)
	case "minute":
		return d.Round(time.Minute)
	case "hour":
		return d.Round(time.Hour)
	default:
		return d.Round(time.Nanosecond)
	}
}

func daySummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	mgr := ctx.CreateManager()
	date := mgr.TimeProvider.Now().Time.In(loc)
	rawDate := strings.TrimSpace(r.PathValue("date"))

	if rawDate != "" {
		date, _ = time.ParseInLocation(time.DateOnly, rawDate, loc)
	}

	durations := map[string]time.Duration{}
	overallDuration := time.Duration(0)
	response := DaySummaryResponse{}

	mgr.ContextStore.Read(func(s *ctx_model.State) error {
		for ctxId, _ := range s.Contexts {
			d, err := mgr.GetIntervalDurationsByDate(s, ctxId, ctx_model.ZonedTime{Time: date, Timezone: loc.String()})
			util.Checkm(err, "Unable to get interval durations for context "+ctxId)
			durations[ctxId] = roundDuration(d, "nanosecond")
		}
		return nil
	})

	sortedIds := make([]string, 0, len(durations))
	for k := range durations {
		sortedIds = append(sortedIds, k)
	}
	sort.Strings(sortedIds)

	for _, c := range sortedIds {
		d := durations[c]
		ctx, _ := mgr.Ctx(c)
		if d > 0 {
			overallDuration += d
			mgr.ContextStore.Read(func(s *ctx_model.State) error {
				response.Contexts = append(response.Contexts, ContextSummary{
					Id:          c,
					Description: ctx.Description,
					Intervals:   mgr.GetIntervalsByDate(s, c, ctx_model.ZonedTime{Time: date, Timezone: loc.String()}),
					Duration:    d,
				})
				return nil
			})
		}
	}

	response.Duration = overallDuration

	json.NewEncoder(w).Encode(response)
}

func Serve() {
	content, err := fs.Sub(staticFiles, "ui/ctx-dashboard/dist/ctx-dashboard")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(content))

	http.Handle("/", spaHandler(content, fs))
	http.HandleFunc("/api/context/list", contextList)
	http.HandleFunc("/api/context/current", currentContext)
	http.HandleFunc("/api/context/free", freeContext)
	http.HandleFunc("/api/context/switch", switchContext)
	http.HandleFunc("/api/context/createAndSwitch", createAndSwitchContext)
	http.HandleFunc("/api/context/interval", updateInterval)
	http.HandleFunc("/api/summary/day/{date}", daySummary)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type SwitchRequest struct {
	Id string `json:"id"`
}

type createAndSwitchRequest struct {
	Description string `json:"description"`
}

type EditIntervalRequest struct {
	Id         string              `json:"contextId"`
	IntervalId string              `json:"intervalId"`
	Start      ctx_model.ZonedTime `json:"start"`
	End        ctx_model.ZonedTime `json:"end"`
}

type ContextSummary struct {
	Id          string               `json:"id"`
	Description string               `json:"description"`
	Intervals   []ctx_model.Interval `json:"intervals"`
	Duration    time.Duration        `json:"duration"`
}

type DaySummaryResponse struct {
	Contexts []ContextSummary `json:"contexts"`
	Duration time.Duration    `json:"duration"`
}

func switchContext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var p SwitchRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	mgr := ctx.CreateManager()
	mgr.Switch(p.Id)
}

func freeContext(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	mgr := ctx.CreateManager()
	mgr.Free()
}

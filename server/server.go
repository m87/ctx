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
	"strconv"
	"strings"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	ctxtime "github.com/m87/ctx/time"
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

	manager := bootstrap.CreateManager()

	manager.WithSession(func(session core.Session) error {
		output := []core.Context{}
		for _, c := range session.State.Contexts {
			output = append(output, c)
		}
		json.NewEncoder(w).Encode(output)
		return nil
	})
}

func currentContext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	manager := bootstrap.CreateManager()

	res := CurrentContextResponse{
		CurrentDuration: 0,
	}
	manager.WithSession(func(session core.Session) error {
		if session.State.CurrentId != "" {
			currentCtx := session.MustGetCtx(session.State.CurrentId)
			res.Context = currentCtx

			for _, interval := range currentCtx.Intervals {
				if interval.End.Time.IsZero() {
					res.CurrentDuration = manager.TimeProvider.Now().Time.Sub(interval.Start.Time)
				}
			}

			json.NewEncoder(w).Encode(res)
		} else {
			json.NewEncoder(w).Encode(nil)
		}
		return nil
	})

}

func createAndSwitchContext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	manager := bootstrap.CreateManager()

	var p createAndSwitchRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	manager.WithArchiveSession(func(session core.Session) error {
		return session.CreateIfNotExistsAndSwitch(util.GenerateId(p.Description), p.Description)
	})
}

func renameContext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	manager := bootstrap.CreateManager()

	var p RenameRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	manager.WithSession(func(session core.Session) error {
		session.RenameContext(p.CtxId, util.GenerateId(p.Name), p.Name)
		return nil
	})

}

func updateInterval(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	manager := bootstrap.CreateManager()

	var p EditIntervalRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	manager.WithSession(func(session core.Session) error {
		session.EditContextInterval(p.Id, p.IntervalId, p.Start, p.End)
		return nil
	})
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

func intervalsByDate(date time.Time) (IntervalsResponseEntry, error) {
	response := IntervalsResponseEntry{}

	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	manager := bootstrap.CreateManager()

	response.Date = date.Format(time.DateOnly)

	manager.WithSession(func(session core.Session) error {
		for ctxId, _ := range session.State.Contexts {
			intervals := session.GetIntervalsByDate(ctxId, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
			for _, i := range intervals {
				response.Intervals = append(response.Intervals, IntervalEntry{
					Id:          i.Id,
					CtxId:       ctxId,
					Description: session.State.Contexts[ctxId].Description,
					Interval:    i,
				})
			}

		}
		return nil
	})

	return response, nil
}

func intervals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	mgr := bootstrap.CreateManager()
	date := mgr.TimeProvider.Now().Time.In(loc)
	rawDate := strings.TrimSpace(r.PathValue("date"))

	if rawDate != "" {
		date, _ = time.ParseInLocation(time.DateOnly, rawDate, loc)
	}

	interval, _ := intervalsByDate(date)
	response := IntervalsResponse{}
	response.Days = append(response.Days, interval)

	json.NewEncoder(w).Encode(response)
}

func daySUmmaryByDate(date time.Time) (DaySummaryResponse, error) {

	manager := bootstrap.CreateManager()
	durations := map[string]time.Duration{}
	overallDuration := time.Duration(0)
	response := DaySummaryResponse{}
	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}

	manager.WithSession(func(session core.Session) error {
		for ctxId, _ := range session.State.Contexts {
			d, err := session.GetIntervalDurationsByDate(ctxId, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
			util.Checkm(err, "Unable to get interval durations for context "+ctxId)
			durations[ctxId] = roundDuration(d, "nanosecond")
		}
		sortedIds := make([]string, 0, len(durations))
		for k := range durations {
			sortedIds = append(sortedIds, k)
		}
		sort.Strings(sortedIds)

		for _, c := range sortedIds {
			d := durations[c]
			ctx := session.MustGetCtx(c)
			if d > 0 {
				overallDuration += d
				intervals := session.GetIntervalsByDate(c, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
				intervalsMap := make(map[string]core.Interval)
				for _, interval := range intervals {
					intervalsMap[interval.Id] = interval
				}
				response.Contexts = append(response.Contexts, core.Context{
					Id:          c,
					Description: ctx.Description,
					Intervals:   intervalsMap,
					Duration:    d,
				})
			}
		}

		response.Duration = overallDuration
		return nil
	})

	return response, nil
}

func dayListSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	mgr := bootstrap.CreateManager()

	response := make(map[string]int)
	mgr.WithSession(func(session core.Session) error {
		response = session.GetContextCountByDateMap()
		return nil
	})

	json.NewEncoder(w).Encode(response)
}

func daySummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	mgr := bootstrap.CreateManager()
	date := mgr.TimeProvider.Now().Time.In(loc)
	rawDate := strings.TrimSpace(r.PathValue("date"))

	if rawDate != "" {
		date, _ = time.ParseInLocation(time.DateOnly, rawDate, loc)
	}

	response, _ := daySUmmaryByDate(date)

	json.NewEncoder(w).Encode(response)
}

func recentIntervals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	mgr := bootstrap.CreateManager()
	date := mgr.TimeProvider.Now().Time.In(loc)
	rawDate := strings.TrimSpace(r.PathValue("date"))

	if rawDate != "" {
		date, _ = time.ParseInLocation(time.DateOnly, rawDate, loc)
	}

	n := 10
	rawN := strings.TrimSpace(r.PathValue("n"))
	if rawN != "" {
		n, _ = strconv.Atoi(rawN)
	}

	response := IntervalsResponse{}

	for i := 0; i < n; i++ {
		d := date.AddDate(0, 0, -i)
		intervals, err := intervalsByDate(d)
		if err != nil {
			http.Error(w, "Error fetching intervals for date "+d.Format(time.DateOnly), http.StatusInternalServerError)
			return
		}
		response.Days = append(response.Days, intervals)
	}
	json.NewEncoder(w).Encode(response)

}

func moveInterval(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var p MoveIntervalRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	manager := bootstrap.CreateManager()
	manager.WithArchiveSession(func(session core.Session) error {
		return session.MoveIntervalById(p.Src, p.Target, p.Id)
	})
}

func splitInterval(w http.ResponseWriter, r *http.Request) {
	var p SplitRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ctxId := strings.TrimSpace(r.PathValue("ctxId"))
	id := strings.TrimSpace(r.PathValue("id"))

	manager := bootstrap.CreateManager()
	manager.WithSession(func(session core.Session) error {
		return session.SplitContextIntervalById(ctxId, id, p.Split.H, p.Split.M, p.Split.S)
	})

}

func deleteInterval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctxId := strings.TrimSpace(r.PathValue("ctxId"))
	id := strings.TrimSpace(r.PathValue("id"))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	bootstrap.CreateManager().WithSession(func(session core.Session) error {
		return session.DeleteInterval(ctxId, id)
	})

}

func version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := VersionResponse{}
	response.Version = core.Version

	json.NewEncoder(w).Encode(response)

}

func Serve(manager *core.ContextManager, port string) {
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
	http.HandleFunc("/api/context/rename", renameContext)
	http.HandleFunc("/api/summary/day/{date}", daySummary)
	http.HandleFunc("/api/summary/day", daySummary)
	http.HandleFunc("/api/summary/day/list", dayListSummary)
	http.HandleFunc("/api/intervals/{date}", intervals)
	http.HandleFunc("/api/intervals", intervals)
	http.HandleFunc("/api/intervals/recent/{n}", recentIntervals)
	http.HandleFunc("/api/intervals/move", moveInterval)
	http.HandleFunc("/api/intervals/{ctxId}/{id}", deleteInterval)
	http.HandleFunc("/api/intervals/{ctxId}/{id}/split", splitInterval)
	http.HandleFunc("/api/version", version)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type Split struct {
	H int `json:"h"`
	M int `json:"m"`
	S int `json:"s"`
}

type SplitRequest struct {
	Split Split `json:"split"`
}

type CurrentContextResponse struct {
	Context         core.Context  `json:"context"`
	CurrentDuration time.Duration `json:"currentDuration"`
}

type VersionResponse struct {
	Version string `json:"version"`
}

type SwitchRequest struct {
	Id string `json:"id"`
}

type MoveIntervalRequest struct {
	Id     string `json:"id"`
	Src    string `json:"src"`
	Target string `json:"target"`
}

type createAndSwitchRequest struct {
	Description string `json:"description"`
}

type EditIntervalRequest struct {
	Id         string            `json:"contextId"`
	IntervalId string            `json:"intervalId"`
	Start      ctxtime.ZonedTime `json:"start"`
	End        ctxtime.ZonedTime `json:"end"`
}

type DaySummaryResponse struct {
	Contexts []core.Context `json:"contexts"`
	Duration time.Duration  `json:"duration"`
}

type DaysSyummaryResponse struct {
	Sumarries map[string]DaySummaryResponse `json:"summaries"`
}

type IntervalsResponse struct {
	Days []IntervalsResponseEntry `json:"days"`
}

type IntervalsResponseEntry struct {
	Date      string          `json:"date"`
	Intervals []IntervalEntry `json:"intervals"`
}

type IntervalEntry struct {
	Id          string        `json:"id"`
	CtxId       string        `json:"ctxId"`
	Description string        `json:"description"`
	Interval    core.Interval `json:"interval"`
}

type RenameRequest struct {
	CtxId string `json:"ctxId"`
	Name  string `json:"name"`
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

	manager := bootstrap.CreateManager()
	manager.WithSession(func(session core.Session) error {
		return session.Switch(p.Id)
	})
}

func freeContext(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	manager := bootstrap.CreateManager()
	manager.WithSession(func(session core.Session) error {
		return session.Free()
	})
}

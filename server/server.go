package server

import (
	"log"
	"net/http"
	"time"

	"github.com/m87/ctx/core"
	ctxtime "github.com/m87/ctx/time"
)

type Server struct {
	mgr *core.ContextManager
	mux *http.ServeMux
}

func New(mgr *core.ContextManager) *Server {
	s := &Server{mgr: mgr, mux: http.NewServeMux()}

	content, fsHandler := mustStaticFS()
	s.mux.Handle("/", spaHandler(content, fsHandler))

	ctxMux := http.NewServeMux()
	registerContext(ctxMux, s.mgr)
	s.mux.Handle("/api/context/", http.StripPrefix("/api/context", ctxMux))

	sumMux := http.NewServeMux()
	registerSummary(sumMux, s.mgr)
	s.mux.Handle("/api/summary/", http.StripPrefix("/api/summary", sumMux))

	intMux := http.NewServeMux()
	registerIntervals(intMux, s.mgr)
	s.mux.Handle("/api/intervals/", http.StripPrefix("/api/intervals", intMux))

	s.mux.HandleFunc("/api/version", s.version)

	return s
}

func (s *Server) Handler() http.Handler {
	var h http.Handler = s.mux
	h = withLogging(h)
	return h
}

func (s *Server) Listen(addr string) error {
	log.Printf("listening on %s", addr)
	return http.ListenAndServe(addr, s.Handler())
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
	Contexts      []core.Context `json:"contexts"`
	OtherContexts []core.Context `json:"otherContexts"`
	Duration      time.Duration  `json:"duration"`
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

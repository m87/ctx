package time

import (
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type TimeProvider interface {
	Now() ZonedTime
}

type ZonedTime struct {
	Time     time.Time `json:"time"`
	Timezone string    `json:"timezone"`
}

func DetectTimezoneName() string {
	switch runtime.GOOS {
	case "linux", "darwin":
		return detectUnixTimezone()
	default:
		return "UTC"
	}
}

func detectUnixTimezone() string {
	if tz := os.Getenv("TZ"); tz != "" {
		return tz
	}

	out, err := exec.Command("readlink", "-f", "/etc/localtime").Output()
	if err != nil {
		return "UTC"
	}

	path := strings.TrimSpace(string(out))
	const zoneinfoPrefix = "/usr/share/zoneinfo/"
	if strings.HasPrefix(path, zoneinfoPrefix) {
		return strings.TrimPrefix(path, zoneinfoPrefix)
	}

	return "UTC"
}

func (zt ZonedTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Time     string `json:"time"`
		Timezone string `json:"timezone"`
	}{
		Time:     zt.Time.Format(time.RFC3339),
		Timezone: zt.Time.Location().String(),
	})
}

func (zt *ZonedTime) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Time     string `json:"time"`
		Timezone string `json:"timezone"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	loc, err := time.LoadLocation(tmp.Timezone)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation(time.RFC3339, tmp.Time, loc)
	if err != nil {
		return err
	}
	zt.Time = t
	zt.Timezone = tmp.Timezone
	return nil
}

type RealTimeProvider struct{}

func (provider *RealTimeProvider) Now() ZonedTime {
	loc, err := time.LoadLocation(DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	return ZonedTime{Time: time.Now().In(loc), Timezone: loc.String()}
}

func NewTimer() *RealTimeProvider {
	return &RealTimeProvider{}
}

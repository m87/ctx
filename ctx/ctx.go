package ctx

import (
	"encoding/json"
	"os"
	"time"

	"github.com/m87/ctx/util"
)



type ContextState int

const (
  ACTIVE ContextState = iota
  FINISHED
)

type Interval struct {
  start time.Time
  end time.Time
  duration int32
}


type Context struct {
  Id string
  Description string
  State ContextState
  Duration int32
  Intervals []Interval
}


type State struct {
  Contexts []Context
  Current Context
}


func Load() State {
  data, err := os.ReadFile("test/state")
  util.Check(err, "Unable to read state file");

  state := State {}
  err = json.Unmarshal(data, &state)
  util.Check(err, "Unable to parse state file");

  return state
}

func Save(state State) {
  data, err := json.Marshal(state)
  if err != nil {
    panic(err)
  }
  os.WriteFile("test/state", data, 0644)
}

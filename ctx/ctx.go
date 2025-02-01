package ctx

import (
	"encoding/json"
	"os"
	"time"

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
  id string
  description string
  state ContextState
  duration int32
  intervals []Interval
}


type State struct {
  contexts []Context
  current Context
}



func Save(state State) {
  data, err := json.Marshal(state)
  if err != nil {
    panic(err)
  }
  os.WriteFile("test/state", data, 0644)

}

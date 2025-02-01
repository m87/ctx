package util

import "log"


func Check(err error, msg string) {
  if err != nil {
    log.Fatal(msg)
  }
}

package main

import (
	"time"
)

type message struct {
	Name string
	Msg  string
	When time.Time
}

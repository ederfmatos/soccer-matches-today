package main

import "time"

func Now() time.Time {
	return time.Now().UTC().Add(time.Hour * -3)
}

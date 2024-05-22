package services

import (
	"time"
)

func GetCreationTime() string {
	return time.Now().Format("2006-01-02T15:04:05.999Z")
}

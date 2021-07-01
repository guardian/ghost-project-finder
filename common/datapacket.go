package common

import "time"

type DataPacket struct {
	Hostname string    `json:"hostname"`
	Fullpath string    `json:"fullpath"`
	Size     int64     `json:"size"`
	Modtime  time.Time `json:"modtime"`
}

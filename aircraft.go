package main

import (
	"fmt"
	"math"
	"time"
)

type aircraftData struct {
	icaoAddr uint32

	callsign string

	eRawLat uint32
	eRawLon uint32
	oRawLat uint32
	oRawLon uint32

	latitude  float64
	longitude float64
	altitude  int32

	lastPing time.Time
	lastPos  time.Time

	mlat bool
}

type aircraftList []*aircraftData
type aircraftMap map[uint32]*aircraftData

func (a aircraftList) Len() int {
	return len(a)
}

func (a aircraftList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a aircraftList) Less(i, j int) bool {
	if a[i].latitude != math.MaxFloat64 && a[j].latitude != math.MaxFloat64 {
		return sortAircraftByDistance(a, i, j)
	} else if a[i].latitude != math.MaxFloat64 && a[j].latitude == math.MaxFloat64 {
		return true
	} else if a[i].latitude == math.MaxFloat64 && a[j].latitude != math.MaxFloat64 {
		return false
	}
	return sortAircraftByCallsign(a, i, j)
}

func sortAircraftByDistance(a aircraftList, i, j int) bool {
	return greatcircle(a[i].latitude, a[i].longitude, *baseLat, *baseLon) <
		greatcircle(a[j].latitude, a[j].longitude, *baseLat, *baseLon)
}

func sortAircraftByCallsign(a aircraftList, i, j int) bool {
	if a[i].callsign != "" && a[j].callsign != "" {
		return a[i].callsign < a[j].callsign
	} else if a[i].callsign != "" && a[j].callsign == "" {
		return true
	} else if a[i].callsign == "" && a[j].callsign != "" {
		return false
	}
	return fmt.Sprintf("%06x", a[i].icaoAddr) < fmt.Sprintf("%06x", a[j].icaoAddr)
}

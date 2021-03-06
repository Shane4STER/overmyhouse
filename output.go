package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

func durationSecondsElapsed(since time.Duration) string {
	sec := uint8(since.Seconds())
	if sec == math.MaxUint8 {
		return "-"
	} else {
		return fmt.Sprintf("%4d", sec)
	}
}

func printOverhead(knownAircraft *aircraftMap) {
	sortedAircraft := make(aircraftList, 0, len(*knownAircraft))

	for _, aircraft := range *knownAircraft {
		sortedAircraft = append(sortedAircraft, aircraft)
	}

	sort.Sort(sortedAircraft)

	for _, aircraft := range sortedAircraft {
		stale := (time.Since(aircraft.lastPos) > time.Duration((10)*time.Second))
		extraStale := (time.Since(aircraft.lastPos) > (time.Duration(20) * time.Second))

		aircraftHasLocation := (aircraft.latitude != math.MaxFloat64 &&
			aircraft.longitude != math.MaxFloat64)
		aircraftHasAltitude := aircraft.altitude != math.MaxInt32

		if aircraftHasLocation && aircraftHasAltitude {
			sLatLon := fmt.Sprintf("%f,%f", aircraft.latitude, aircraft.longitude)
			sAlt := fmt.Sprintf("%d", aircraft.altitude)

			distance := greatcircle(aircraft.latitude, aircraft.longitude,
				*baseLat, *baseLon)

			tPos := time.Since(aircraft.lastPos)

			if !stale && !extraStale && metersInMiles(distance) < 6 {
				fmt.Printf("%06x\t%8s\t%s%s\t%s\t%3.2f\t%s\n",
					aircraft.icaoAddr, aircraft.callsign,
					sLatLon, sAlt, metersInMiles(distance),
					durationSecondsElapsed(tPos))
			}
		}
	}
}

func printAircraftTable(knownAircraft *aircraftMap) {
	fmt.Print("\x1b[H\x1b[2J")
	fmt.Println("ICAO \tCallsign\tLocation\t\tAlt\tDistance   Time")

	sortedAircraft := make(aircraftList, 0, len(*knownAircraft))

	for _, aircraft := range *knownAircraft {
		sortedAircraft = append(sortedAircraft, aircraft)
	}

	sort.Sort(sortedAircraft)

	for _, aircraft := range sortedAircraft {
		stale := (time.Since(aircraft.lastPos) > time.Duration((10)*time.Second))
		extraStale := (time.Since(aircraft.lastPos) > (time.Duration(20) * time.Second))

		aircraftHasLocation := (aircraft.latitude != math.MaxFloat64 &&
			aircraft.longitude != math.MaxFloat64)
		aircraftHasAltitude := aircraft.altitude != math.MaxInt32

		if aircraft.callsign != "" || aircraftHasLocation || aircraftHasAltitude {
			var sLatLon string
			var sAlt string

			if aircraftHasLocation {
				sLatLon = fmt.Sprintf("%f,%f", aircraft.latitude, aircraft.longitude)
			} else {
				sLatLon = "---.------, ---.------"
			}
			if aircraftHasAltitude {
				sAlt = fmt.Sprintf("%d", aircraft.altitude)
			} else {
				sAlt = "-----"
			}

			distance := greatcircle(aircraft.latitude, aircraft.longitude,
				*baseLat, *baseLon)

			isMlat := ""
			if aircraft.mlat {
				isMlat = "^"
			}

			tPos := time.Since(aircraft.lastPos)

			if !stale && !extraStale {
				fmt.Printf("%06x\t%8s\t%s%s\t%s\t%3.2f\t%s\n",
					aircraft.icaoAddr, aircraft.callsign,
					sLatLon, isMlat, sAlt, metersInMiles(distance),
					durationSecondsElapsed(tPos))
			} else if stale && !extraStale {
				fmt.Printf("%06x\t%8s\t%s%s?\t%s\t%3.2f?\t%s\n",
					aircraft.icaoAddr, aircraft.callsign,
					sLatLon, isMlat, sAlt, metersInMiles(distance),
					durationSecondsElapsed(tPos))
			} else if extraStale {
				fmt.Printf("%06x\t%8s\t%s%s?\t%s\t%3.2f?\t%s…\n",
					aircraft.icaoAddr, aircraft.callsign,
					sLatLon, isMlat, sAlt, metersInMiles(distance),
					durationSecondsElapsed(tPos))
			}
		}
	}
}

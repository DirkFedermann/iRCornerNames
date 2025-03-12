package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mpapenbr/goirsdk/irsdk"
	"gopkg.in/yaml.v3"
)

func connectToIracing() bool {
	logDebug("[iracing] [connectToIracing] Try connecting to iRacing")
	client := http.Client{Timeout: time.Second * 10}
	if simRunning, err := irsdk.IsSimRunning(context.Background(), &client); err != nil {
		panic(err)
	} else if simRunning {
		logDebug("[iracing] [connectToIracing] iRacing is running and connected")
		return true
	}
	logDebug("[iracing] [connectToIracing] iRacing is not running")
	return false
}

func getCamCarIdx(iRApi *irsdk.Irsdk) int32 {
	logDebug("[iracing] [getCamCarIdx] Getting CamCarIdx")
	camCarIdx, err := iRApi.GetValue("CamCarIdx")
	if err != nil {
		logDebug(fmt.Sprintf("[iracing] [getCamCarIdx] Error getting CamCarIdx: %v", err))
		return -1
	}
	logDebug("[iracing] [getCamCarIdx] CamCarIdx: %v", camCarIdx)
	return camCarIdx.(int32)
}

func getCarIdxLapDistPct(iRApi *irsdk.Irsdk, carIdx int32) []float32 {
	logDebug("[iracing] [getCarIdxLapDistPct] Getting CarIdxLapDistPct for carIdx: %v", carIdx)
	carIdxLapDistPct, err := iRApi.GetValue("CarIdxLapDistPct")
	if err != nil {
		logDebug(fmt.Sprintf("[iracing] [getCarIdxLapDistPct] Error getting CarIdxLapDistPct: %v", err))
		return nil
	}
	return []float32{carIdxLapDistPct.([]float32)[carIdx]}
}

// func getLapDist(iRApi *irsdk.Irsdk) float32 {
// 	logDebug("[iracing] [getLapDist] Getting LapDist")
// 	lapDist, err := iRApi.GetValue("LapDist")
// 	if err != nil {
// 		logDebug(fmt.Sprintf("[iracing] [getLapDist] Error getting LapDist: %v", err))
// 		return -1
// 	}
// 	logDebug("[iracing] [getLapDist] LapDist: %v", lapDist)
// 	return lapDist.(float32)
// }

// func getLapDistPct(iRApi *irsdk.Irsdk) float32 {
// 	logDebug("[iracing] [getLapDistPct] Getting LapDistPct")
// 	lapDistPct, err := iRApi.GetValue("LapDistPct")
// 	if err != nil {
// 		logDebug(fmt.Sprintf("[iracing] [getLapDistPct] Error getting LapDistPct: %v", err))
// 		return -1
// 	}
// 	logDebug("[iracing] [getLapDistPct] LapDistPct: %v", lapDistPct)
// 	return lapDistPct.(float32)
// }

type WeekendInfo struct {
	WeekendInfo struct {
		TrackID     int    `yaml:"TrackID"`
		TrackLength string `yaml:"TrackLength"`
	} `yaml:"WeekendInfo"`
}

func getWeekendInfo(iRApi *irsdk.Irsdk) WeekendInfo {

	logDebug("[iracing] [getWeekendInfo] Getting weekend info")

	yamlData := iRApi.GetYamlString()

	var weekendInfo WeekendInfo
	data := yaml.NewDecoder(strings.NewReader(yamlData))
	err := data.Decode(&weekendInfo)
	if err != nil {
		logDebug("[iracing] [getWeekendInfo] Error decoding YAML data: %v", err)
		return WeekendInfo{}
	}
	logDebug("[iracing] [getWeekendInfo] Decoded YAML data: %+v", weekendInfo)

	return weekendInfo
}

func getTrackId(iRApi *irsdk.Irsdk) int {
	weekendInfo := getWeekendInfo(iRApi)
	logDebug("[iracing] [getTrackId] TrackID: %v", weekendInfo.WeekendInfo.TrackID)
	return weekendInfo.WeekendInfo.TrackID
}

func getTrackLength(iRApi *irsdk.Irsdk) int {
	weekendInfo := getWeekendInfo(iRApi)
	trackLengthString := weekendInfo.WeekendInfo.TrackLength
	trackLengthParts := strings.Split(trackLengthString, " ")
	var trackLengthMeters int
	if len(trackLengthParts) == 2 {
		trackLengthValue, err := strconv.ParseFloat(trackLengthParts[0], 64)
		if err == nil {
			trackLengthMeters = int(trackLengthValue * 1000)
		}
	}

	return trackLengthMeters
}

func getCornerName(csvrecords [][]string, calculatedLapDist float32) string {
	for i, record := range csvrecords {
		if i == 0 {
			continue
		}
		start, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			logDebug("[iracing] [getCornerName] Error parsing start value: %v", err)
			continue
		}
		end, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			logDebug("[iracing] [getCornerName] Error parsing end value: %v", err)
			continue
		}
		if calculatedLapDist >= float32(start) && calculatedLapDist <= float32(end) {
			logDebug("[iracing] [getCornerName] Corner name: %s", record[0])
			return record[0]
		}
	}
	return ""
}

func getCalculatedLapDist(tracklength int, carIdxLapDistPct float32) float32 {
	return float32(tracklength) * float32(carIdxLapDistPct) / 100
}

func getPitStatus(iRApi *irsdk.Irsdk, carIdx int) bool {
	value, err := iRApi.GetValue("CarIdxOnPitRoad")
	pitStatus := value.([]bool)[carIdx]
	if err != nil {
		logDebug("[iracing] [getPitStatus] Error getting pit status: %v", err)
		return false
	}
	return pitStatus
}

// func getSessionData(iRApi *irsdk.Irsdk) {
// 	logDebug("[iracing] [getSessionData] Getting session data")

// 	// Get the current watched carIdx
// 	var camCarIdx int32 = getCamCarIdx(iRApi)
// 	logDebug(fmt.Sprintf("[iracing] [getSessionData] CamCarIdx: %v", camCarIdx))

// 	// Get the lap distance percentage for the watched carIdx
// 	var carIdxLapDistPct []float32 = getCarIdxLapDistPct(iRApi, camCarIdx)
// 	logDebug(fmt.Sprintf("[iracing] [getSessionData] CarIdxLapDistPct: %v", carIdxLapDistPct))

// 	// Get the lap distance
// 	var lapDist float32 = getLapDist(iRApi)
// 	logDebug(fmt.Sprintf("[iracing] [getSessionData] LapDist: %v", lapDist))

// 	// Get the lap distance percentage
// 	var lapDistPct float32 = getLapDistPct(iRApi)
// 	logDebug(fmt.Sprintf("[iracing] [getSessionData] LapDistPct: %v", lapDistPct))

// 	// Get the track ID
// 	var trackId int = getTrackId(iRApi)
// 	logDebug(fmt.Sprintf("[iracing] [getSessionData] TrackID: %v", trackId))

// 	// Get the track length
// 	var trackLength int = getTrackLength(iRApi)
// 	logDebug(fmt.Sprintf("[iracing] [getSessionData] TrackLength in meters: %v", trackLength))
// }

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mpapenbr/goirsdk/irsdk"
)

var (
	CalculatedLapDist float32
	Csvrecords        [][]string
	debug             bool = true
	webdebug          bool = true
	showlapdist       bool = true
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
	flag.BoolVar(&webdebug, "webdebug", false, "Enable web debug output")
	flag.BoolVar(&showlapdist, "showlapdist", false, "Show lap distance")
}

func logDebug(message any, args ...any) {
	if debug {
		fmt.Printf("[DEBUG] %v\n", fmt.Sprintf(message.(string), args...))
	}
}

func logWebDebug(message any, args ...any) {
	if webdebug {
		fmt.Printf("[WEBDEBUG] %v\n", fmt.Sprintf(message.(string), args...))
	}
}

func loadCsv(csvfile string) [][]string {
	// load csv
	// csvfile = "tracks/262.csv"
	logDebug("[main] Loading csv file: %s", csvfile)
	file, err := os.Open(csvfile)
	if err != nil {
		logDebug("[main] Error opening file: %v", err)
		fmt.Printf("[main] Error opening file: %v\n", err)
		fmt.Println("Press 'Enter' to exit...")
		fmt.Scanln()
		os.Exit(1)
	} else {
		logDebug("[main] File opened")
	}
	defer file.Close()
	logDebug("[main] File closed")

	// read csv
	logDebug("[main] Reading csv file")
	reader := csv.NewReader(file)
	Csvrecords, err := reader.ReadAll()
	if err != nil {
		logDebug("[main] Error reading file: %v", err)
		return [][]string{}
	} else {
		logDebug("[main] File read")
	}
	return Csvrecords
}

func main() {
	flag.Parse()
	logDebug("[main] Debug output enabled")
	logWebDebug("[main] Web debug output enabled")

	if !irsdk.CheckIfSimIsRunning() {
		fmt.Println("Error: iRacing is not running or the SDK is not connected.")
		os.Exit(1)
	}

	fmt.Println("sadfasdf")
	iRApi := irsdk.NewIrsdk()
	if iRApi == nil {
		fmt.Println("bli")
	} else {
		fmt.Println("blo")
	}
	iRApi.GetData()

	logDebug("[main] Starting web server")
	go webServer()
	logDebug("[main] Web server started")

	// Connect to iRacing
	retries := 1
	maxRetires := 10
	logDebug("[main] Connecting to iRacing")
	for retries <= maxRetires {
		if connectToIracing() {
			logDebug("[main] Connected to iRacing")
			break
		}
		logDebug(fmt.Sprintf("[main] iRacing not running, retrying in %d seconds", retries))
		time.Sleep(time.Second * time.Duration(retries))
		retries++
	}
	if retries > maxRetires {
		logDebug("[main] Failed to connect to iRacing")
	}

	trackId := strconv.Itoa(getTrackId(iRApi))
	logDebug("[main] loading csv file: %s.csv", trackId)
	csvFile := "tracks/" + trackId + ".csv"
	loadCsv(csvFile)

	var camCarIdx int32 = getCamCarIdx(iRApi)
	logDebug("[main] CamCarIdx: %v", camCarIdx)
	var carIdxLapDistPct float32 = getCarIdxLapDistPct(iRApi, camCarIdx)[0] * 100
	logDebug("[main] CarIdxLapDistPct: %v", carIdxLapDistPct)
	var trackLength int = getTrackLength(iRApi)
	logDebug("[main] TrackLength: %v", trackLength)
	calc := float32(trackLength) * carIdxLapDistPct / 100
	logDebug("[main] Calculated meters after S/F: %v", calc)

	// Loop as long as program is running
	for {
		iRApi.GetData()
		var camCarIdx int32 = getCamCarIdx(iRApi)
		logDebug("[main] CamCarIdx: %v", camCarIdx)
		var carIdxLapDistPct float32 = getCarIdxLapDistPct(iRApi, camCarIdx)[0] * 100
		logDebug("[main] CarIdxLapDistPct: %v", carIdxLapDistPct)
		var trackLength int = getTrackLength(iRApi)
		logDebug("[main] TrackLength: %v", trackLength)
		CalculatedLapDist := getCalculatedLapDist(int(trackLength), carIdxLapDistPct)
		logDebug("[main] Calculated meters after S/F: %v", CalculatedLapDist)

		cornerName := getCornerName(Csvrecords, CalculatedLapDist)
		logDebug("[main] Corner name: %s", cornerName)
		time.Sleep(time.Second * 100000)
	}
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/mpapenbr/goirsdk/irsdk"
)

var (
	port int = 8080
)

func init() {
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	logWebDebug("[web] [init] Port: %v", port)
}

func webServer() {
	flag.Parse()
	var iRApi *irsdk.Irsdk

	retries := 1
	maxRetries := 10
	fmt.Println("Waiting for iRacing to start")
	for {
		logWebDebug("[web] [webServer] Trying to connect to iRacing")
		if irsdk.CheckIfSimIsRunning() {
			logWebDebug("[web] [webServer] Connected to iRacing")
			break
		}
		time.Sleep(time.Second * time.Duration(retries))
		if retries <= maxRetries {
			retries++
		}
	}
	iRApi = irsdk.NewIrsdk()
	fmt.Println("iRacing started and connected")

	iRApi.GetData()

	logWebDebug("Reading in trackId and csv")
	trackId := strconv.Itoa(getTrackId(iRApi))
	csvRecords := loadCsv("tracks/" + trackId + ".csv")

	logWebDebug("[web] [webServer] Starting web server")
	fmt.Printf("Starting web server on port %d\n", port)
	fmt.Printf("Go to http://localhost:%d\n", port)

	logWebDebug("[web] [webServer] Starting web server on port %d", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("html/index.html")
		if err != nil {
			http.Error(w, "Error parsing template file", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		updateHandler(w, r, iRApi, csvRecords)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func updateHandler(w http.ResponseWriter, _ *http.Request, iRApi *irsdk.Irsdk, csvRecords [][]string) {
	logWebDebug("[web] [updateHandler] Update request received")

	iRApiGetData := iRApi.GetData()

	logWebDebug("[web] [updateHandler] See if iRacing is still running")
	if !irsdk.CheckIfSimIsRunning() {
		logWebDebug("[web] [updateHandler] iRacing is not running")
		tmpl, err := template.ParseFiles("html/offline.html")
		if err != nil {
			http.Error(w, "Error parsing template file", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, "iRacing is not running")
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
		return
	}

	logWebDebug("[web] [updateHandler] SimIsRunning: %v", iRApiGetData)
	if !iRApiGetData {
		tmpl, err := template.ParseFiles("html/offline.html")
		if err != nil {
			http.Error(w, "Error parsing template file", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, "iRacing is not running")
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
		return
	}

	camCarIdx := getCamCarIdx(iRApi)
	logWebDebug("[web] [updateHandler] CamCarIdx: %d", camCarIdx)
	carIdxLapDistPct := getCarIdxLapDistPct(iRApi, camCarIdx)[0] * 100
	logWebDebug("[web] [updateHandler] CarIdxLapDistPct: %f", carIdxLapDistPct)
	trackLength := int(getTrackLength(iRApi))
	logWebDebug("[web] [updateHandler] Track length: %d", trackLength)
	CalculatedLapDist := getCalculatedLapDist(int(trackLength), carIdxLapDistPct)
	logWebDebug("[web] [updateHandler] Calculated lap distance: %f", CalculatedLapDist)

	cornerName := ""
	logWebDebug("[web] [updateHandler] PitStatus: %v", getPitStatus(iRApi, int(camCarIdx)))
	if getPitStatus(iRApi, int(camCarIdx)) {
		cornerName = "Pits"
	} else {
		cornerName = getCornerName(csvRecords, CalculatedLapDist)
	}
	if showlapdist {
		cornerName = cornerName + " - LapDist: " + fmt.Sprintf("%.0f", getCalculatedLapDist(int(trackLength), carIdxLapDistPct))
	}
	tmpl, err := template.ParseFiles("html/update.html")
	if err != nil {
		http.Error(w, "Error parsing template file", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, cornerName)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
	logWebDebug("[web] [updateHandler] Corner name: %s", cornerName)
}

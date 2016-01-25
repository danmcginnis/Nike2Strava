//Handles moving run data from Nike+ into Strava
package nike2strava

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
    "log"
)

const baseURL = "https://api.nike.com/v1/me/sport/activities"

type nikeDataSimple struct {
	Data []struct {
		ActivityID    string
		ActivityType  string
		StartTime     string
		Status        string
		MetricSummary struct {
			Distance string
			Duration string
		}
	}
	Paging struct {
		Next     string
		Previous string
	}
}

type nikeDataComplete struct {
	Links []struct {
		Rel  string
		Href string
	}
	ActivityID       string
	ActivityType     string
	StartTime        string
	ActivityTimeZone string
	Status           string
	DeviceType       string
	MetricSummary    struct {
		Calories string
		Distance string
		Duration string
	}
	Paging struct {
		Next     string
		Previous string
	}
	Tags []struct {
		TagType  string
		TagValue string
	}
	Metrics []struct {
		IntervalMetric int16
		IntervalUnit   string
		MetricType     string
		Values         []string
	}
	IsGPSActivity  bool
	ElevationLoss  float64
	ElevationGain  float64
	ElevationMax   float64
	ElevationMin   float64
	IntervalMetric float64
	IntervalUnit   string
	Waypoints      []struct {
		Latitude  float64
		Longitude float64
		Elevation float64
	}
}

func makeActivityURL(token string, count int) string {
	return fmt.Sprintf("%s?access_token=%s&count=%d", baseURL, token, count)
}

func makeDetailsURL(token string, activityID string) string {
	return fmt.Sprintf("%s/%s?access_token=%s", baseURL, activityID, token)
}

func makeGpsURL(token string, activityID string) string {
	return fmt.Sprintf("%s/%s/gps?access_token=%s", baseURL, activityID, token)
}

func findInterval(metric float64, interval string) time.Duration {
	switch interval {
	case "SEC":
		return time.Duration(metric) * time.Second
	case "MIN":
		return time.Duration(metric) * time.Minute
	default:
		//assume the most logical GPS measurement interval
		return time.Duration(metric) * time.Second
	}
}

func formatTimeGPS(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func (nikeActs nikeDataComplete) makeGPX() {
	//todo: write to the screen while we develop, eventually write to a file with
	//  the activity number as the file name
	const longForm = "2006-01-02T15:04:05Z"
	//Go time expects the refence time to be Mon Jan 2 15:04:05 MST 2006
	t, _ := time.Parse(longForm, nikeActs.StartTime)
	s := findInterval(nikeActs.IntervalMetric, nikeActs.IntervalUnit)

	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Println(`<gpx creator="Nike2Strava" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3">`)
	fmt.Println(" <metadata>")
	fmt.Println("  <time>" + nikeActs.StartTime + "</time>")
	fmt.Println(" </metadata>")
	fmt.Println(" <trk>")
	fmt.Println("  <name>" + nikeActs.StartTime + "</name>")
	fmt.Println("  <trkseg>")
	for _, m := range nikeActs.Waypoints {
		fmt.Println("   <trkpt lat=" + `"` + strconv.FormatFloat(m.Latitude, 'f', 7, 64) + `"` + " lon=" + `"` + strconv.FormatFloat(m.Longitude, 'f', 7, 64) + `">`)
		fmt.Println("    <ele>" + strconv.FormatFloat(m.Elevation, 'f', 1, 64) + "</ele>")
		fmt.Println("    <time>" + formatTimeGPS(t) + "</time>")
		fmt.Println("   </trkpt>")
		t = t.Add(s)
	}
	fmt.Println("  </trkseg>")
	fmt.Println(" </trk>")
	fmt.Println("</gpx>")
}

func (nikeActs nikeDataComplete) printDetails() {
	fmt.Println("Activity ID:", nikeActs.ActivityID)
	fmt.Println("Distance:", nikeActs.MetricSummary.Distance, "km")
	fmt.Println("Date:", nikeActs.StartTime)
	for _, m := range nikeActs.Tags {
		fmt.Println(strings.ToLower(m.TagType), ":", strings.ToLower(m.TagValue))
	}
	if nikeActs.IsGPSActivity {
		fmt.Println("Activity has GPS data:", nikeActs.IsGPSActivity)
		fmt.Println("GPS read interval is", nikeActs.IntervalMetric, strings.ToLower(nikeActs.IntervalUnit))
		fmt.Println(nikeActs.Waypoints[0].Latitude, nikeActs.Waypoints[0].Longitude)
		fmt.Println(nikeActs.Waypoints[1].Latitude, nikeActs.Waypoints[1].Longitude)
	}
	fmt.Println()
}

func getDetails(token string, activityID string, debug bool) {
	var nikeActs nikeDataComplete

	if debug {
		json.Unmarshal(NikeActDetails, &nikeActs)
		json.Unmarshal(NikeActGPS, &nikeActs)
	} else {
		url := makeDetailsURL(token, activityID)
		res, err := http.Get(url)
		if err != nil {
			panic(err.Error())
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err.Error())
		}
		json.Unmarshal(body, &nikeActs)

		if nikeActs.IsGPSActivity {
			gpsURL := makeGpsURL(token, activityID)
			res, err = http.Get(gpsURL)
			if err != nil {
				panic(err.Error())
			}
			body, err = ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err.Error())
			}
			json.Unmarshal(body, &nikeActs)
		}
	}
	if !debug {
		nikeActs.printDetails()
	}
	nikeActs.makeGPX()
}

func wrangleJSON(token string, numRecords int, debug bool) {

	var nikeList nikeDataSimple

	if debug {
		json.Unmarshal(NikeBasic1, &nikeList)
		log.Print("Using Local Test Data")
	} else {
		log.Print("\nUsing live data from Nike+ API\n")
		url := makeActivityURL(token, numRecords)
		res, err := http.Get(url)

		if err != nil {
			panic(err.Error())
		}
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			panic(err.Error())
		}

		json.Unmarshal(body, &nikeList)
	}

	for _, m := range nikeList.Data {
		if m.ActivityType == "RUN" {
			getDetails(token, m.ActivityID, debug)
		}
	}
}

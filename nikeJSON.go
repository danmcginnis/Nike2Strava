package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const baseURL = "https://api.nike.com/v1/me/sport/activities"

//var loginURL = "https://developer.nike.com/content/nike-developer-cq/us/en_us/index/login.html"

type nikeDataSimple struct {
	Data []struct {
		ActivityId    string
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
	ActivityId       string
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
		IntervalMetric float64
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

/*func main() {
	var nike nikeDataSimple
	json.Unmarshal([]byte(NikeBasic10), &nike)
	fmt.Println(nike)
}*/

func makeActivityURL(token string, count int) string {
	return baseURL + "?access_token=" + token + "&count=" + strconv.Itoa(count)
}

func makeDetailsURL(token string, activityID string) string {
	return baseURL + "/" + activityID + "?access_token=" + token
}

func makeGpsURL(token string, activityID string) string {
	return baseURL + "/" + activityID + "/gps?access_token=" + token
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

	fmt.Println("Activity ID:", nikeActs.ActivityId)
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

func wrangleJSON(token string, numRecords int, debug bool) {

	var nikeList nikeDataSimple

	if debug {
		json.Unmarshal(NikeBasic5, &nikeList)
	} else {
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
			getDetails(token, m.ActivityId, debug)
		}
	}
}

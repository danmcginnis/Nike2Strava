package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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

func wrangleJSON(token string) {
  
	var nikeList nikeDataSimple
	var ActList = "https://api.nike.com/v1/me/sport/activities?access_token=ShQeydt0MoWobetAG39mWX26Yhss&count=20"
 	var ActDetails = "https://api.nike.com/v1/me/sport/activities/"

	url := ActList
	res, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(body, &nikeList)
	for _, m := range nikeList.Data {
		if m.ActivityType == "RUN" {
			var nikeActs nikeDataComplete

			url := ActDetails + m.ActivityId + "?access_token=" + token
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
				gpsUrl := ActDetails + m.ActivityId + "/gps?access_token=" + token
				res, err = http.Get(gpsUrl)

				if err != nil {
					panic(err.Error())
				}
				body, err = ioutil.ReadAll(res.Body)

				if err != nil {
					panic(err.Error())
				}
				json.Unmarshal(body, &nikeActs)
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
			}
			fmt.Println()
		}
	}
}

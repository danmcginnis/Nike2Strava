package main

import (
	"encoding/json"
	"fmt"
  "strings"
)

type nikeDataSimple struct {
  Data []struct {
    ActivityId string
    ActivityType string
    StartTime string
    Status string
    MetricSummary struct {
      Distance string
      Duration string
    }
  }
  Paging struct {
    Next string
    Previous string
  }
}

type nikeDataComplete struct {
  Data []struct {
    Links []struct {
      Rel string
      Href string
    }
    ActivityId string
    ActivityType string
    StartTime string
    ActivityTimeZone string
    Status string
    DeviceType string
    MetricSummary struct {
      Calories string
      Distance string
      Duration string
    }
  }
  Paging struct {
    Next string
    Previous string
  }
  Tags []struct {
    TagType string
    TagValue string
  }
  Metrics []struct {
    IntervalMetric float64
    IntervalUnit string
    MetricType string
    Values []string
  }
  IsGPSActivity bool
  ElevationLoss float64
  ElevationGain float64
  ElevationMax float64
  ElevationMin float64
  IntervalMetric float64
  IntervalUnit string
  Waypoints []struct {
    Latitude float64
    Longitude float64
    Elevation float64
  }
}

func wrangleJSON() {
	var nike nikeDataComplete
  json.Unmarshal(NikeBasic1, &nike)
  json.Unmarshal(NikeActDetails, &nike)
  json.Unmarshal(NikeActGPS, &nike)
  fmt.Println(nike.Data[0].ActivityId)
  fmt.Println("Activity has GPS data:", nike.IsGPSActivity)
  fmt.Println("GPS read interval is", nike.IntervalMetric, strings.ToLower(nike.IntervalUnit))
  for _, m := range nike.Tags {
    fmt.Println(strings.ToLower(m.TagType), strings.ToLower(m.TagValue))
  }
}


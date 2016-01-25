//Handles moving run data from Nike+ into Strava
package nike2strava

import "testing"

func TestMakeActivityURL(t *testing.T) {
    expected := "https://api.nike.com/v1/me/sport/activities?access_token=123456789&count=5"
    actual := makeActivityURL("123456789", 5)
    if actual != expected {
        t.Errorf("Test Failed, expected '%s', got '%s'", expected, actual)
    }
}


func TestMakeDetailsURL(t *testing.T) {
    expected := "https://api.nike.com/v1/me/sport/activities/987654321?access_token=123456789"
    actual := makeDetailsURL("123456789", "987654321")
    if actual != expected {
        t.Errorf("Test Failed, expected '%s', got '%s'", expected, actual)
    }
}

func TestMakeGpsURL(t *testing.T) {
    expected := "https://api.nike.com/v1/me/sport/activities/987654321/gps?access_token=123456789"
    actual := makeGpsURL("123456789", "987654321")
    if actual != expected {
        t.Errorf("Test Failed, expected '%s', got '%s'", expected, actual)
    }
}
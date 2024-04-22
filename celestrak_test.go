package satellite

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestNewFromOrbitMeanElementsMessage(t *testing.T) {
	var omm CelestrakOrbitMeanElementsMessage
	now := time.Now().UTC()
	rawOmm := `{"OBJECT_NAME":"ISS (ZARYA)","OBJECT_ID":"1998-067A","EPOCH":"2024-04-22T14:33:45.765216","MEAN_MOTION":15.50661244,"ECCENTRICITY":0.0004788,"INCLINATION":51.6403,"RA_OF_ASC_NODE":230.9991,"ARG_OF_PERICENTER":88.2778,"MEAN_ANOMALY":7.6998,"EPHEMERIS_TYPE":0,"CLASSIFICATION_TYPE":"U","NORAD_CAT_ID":25544,"ELEMENT_SET_NO":999,"REV_AT_EPOCH":44992,"BSTAR":0.00060741,"MEAN_MOTION_DOT":0.00035196,"MEAN_MOTION_DDOT":0}`

	if err := json.Unmarshal([]byte(rawOmm), &omm); err != nil {
		t.Errorf("Failed to unmarshal raw OMM: %s", err.Error())
	}

	sat, err := NewFromOrbitMeanElementsMessage(&omm, GravityWGS84)

	if err != nil {
		t.Errorf("Failed create Satellite: %s", err.Error())
	}

	position, _ := Propagate(*sat, now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	gst := GSTimeFromDate(now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	alt, _, latlon := ECIToLLA(position, gst)
	latlondeg := LatLongDeg(latlon)

	fmt.Printf("lat: %.4f, lon %.4f alt %.4f\n", latlondeg.Latitude, latlondeg.Longitude, alt)
}

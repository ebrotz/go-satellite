package satellite

import (
	"time"
)

// Orbit Mean-Elements Message from Celestrak
type CelestrakOrbitMeanElementsMessage struct {
	ObjectName   string  `json:"OBJECT_NAME"`
	ObjectId     string  `json:"OBJECT_ID"`
	Epoch        string  `json:"EPOCH"`
	MeanMotion   float64 `json:"MEAN_MOTION"`
	Eccentricity float64 `json:"ECCENTRICITY"`
	Inclination  float64 `json:"INCLINATION"`
	// Right ascension of the ascending node
	Raan float64 `json:"RA_OF_ASC_NODE"`
	// Argument of periapsis/pericenter
	Periapsis               float64 `json:"ARG_OF_PERICENTER"`
	MeanAnomaly             float64 `json:"MEAN_ANOMALY"`
	EphemerisType           int32   `json:"EPHEMERIS_TYPE"`
	Classification          string  `json:"CLASSIFICATION_TYPE"`
	CatalogId               int64   `json:"NORAD_CAT_ID"`
	ElsetNum                int32   `json:"ELEMENT_SET_NO"`
	RevolutionNumberAtEpoch int32   `json:"REV_AT_EPOCH"`
	// Aerodynamic drag
	Bstar float64 `json:"BSTAR"`
	// First derivative of mean motion
	BallisticCoeifficient float64 `json:"MEAN_MOTION_DOT"`
	// Second derivative of mean motion, area times coefficient of solar radion pressure over mass
	MeanMotionSecondDerivative float64 `json:"MEAN_MOTION_DDOT"`
}

// Creates a new Satellite from an Orbit Mean Elements Message from Celestrak
func NewFromOrbitMeanElementsMessage(omm *CelestrakOrbitMeanElementsMessage, gravConst Gravity) (*Satellite, error) {
	var sat Satellite
	opsmode := "i"
	epoch, err := time.Parse("2006-01-02T15:04:05.999999999", omm.Epoch)

	if err != nil {
		return nil, err
	}

	sat.whichconst = getGravConst(gravConst)

	sat.no = omm.MeanMotion / XPDOTP
	sat.ndot = omm.BallisticCoeifficient / (XPDOTP * 1440.0)
	sat.nddot = omm.MeanMotionSecondDerivative / (XPDOTP * 1440.0 * 1440)

	sat.inclo = omm.Inclination * DEG2RAD
	sat.nodeo = omm.Raan * DEG2RAD
	sat.argpo = omm.Periapsis * DEG2RAD
	sat.mo = omm.MeanAnomaly * DEG2RAD

	sat.jdsatepoch = JDay(epoch.Year(), int(epoch.Month()), epoch.Day(), epoch.Hour(), epoch.Minute(), epoch.Second())

	sgp4init(&opsmode, sat.jdsatepoch-2433281.5, &sat)
	return &sat, nil
}

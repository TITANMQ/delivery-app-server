package utils

import "math"

var earthRadius float32 = 3960.0
var PI float32 = 3.14159265
var degreesToRadians float32 = PI / 180.0
var radiansToDegrees float32 = 180.0 / PI

func MilesToLat(miles float32) float32 {
	return float32((miles / earthRadius) * radiansToDegrees)
}

func MilesToLng(lat float32, miles float32) float32 {
	radius := earthRadius * float32(math.Cos(float64(lat*degreesToRadians)))
	return float32((miles / radius) * radiansToDegrees)
}

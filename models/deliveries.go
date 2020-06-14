package models

import (
	u "backend/utils"
	"fmt"
	"os"
	"strings"

	"github.com/rubenv/opencagedata"
)

//Delivery struct stores delivery data from json to be used to for the database
type Delivery struct {
	DeliveryID        uint    `json:"deliveryID" gorm:"column:deliveryid; auto_increment"`
	ProfileID         int     `json:"profileID" gorm:"column:userid;"`
	DeliveryType      string  `json:"deliveryType" gorm:"column:deliverytype;"`
	CollectionAddress string  `json:"collectionAddress" gorm:"column:collectionaddress;"`
	DeliveryAddress   string  `json:"deliveryAddress" gorm:"column:deliveryaddress;"`
	ExpiryDate        string  `json:"expiryDate" gorm:"column:expirydate;type:timestamp"`
	ExtraDetails      string  `json:"extraDetails" gorm:"column:extradetails;"`
	Accepted          bool    `json:"accepted" gorm:"column:accepted;"`
	Status            string  `json:"status" gorm:"column:status;"`
	CollectionLon     float32 `json:"collectionLon" gorm:"column:collectionlon;"`
	CollectionLat     float32 `json:"collectionLat" gorm:"column:collectionlat;"`
	DeliveryLon       float32 `json:"deliveryLon" gorm:"column:deliverylon;"`
	DeliveryLat       float32 `json:"deliveryLat" gorm:"column:deliverylat;"`
}

//Location struct stores location data
type Location struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

//Journey struct stores journey data
type Journey struct {
	JourneyStart string `json:"journeyStart"`
	JourneyEnd   string `json:"journeyEnd"`
}

//Accepted struct stores accepted_delivery data from a json to be used to for the database
type Accepted struct {
	DeliveryID       uint   `json:"deliveryID" gorm:"column:deliveryid"`
	UserProfileID    uint   `json:"userAccountID" gorm:"-"`
	CourierAccountID uint   `json:"courierAccountID" gorm:"column:profileid"`
	DeliveryDate     string `json:"deliveryDate" gorm:"column:deliverydate;type:timestamp"`
}

// Create adds a new delivery to the database
func (delivery *Delivery) Create() map[string]interface{} {

	lat, lng, err := setLatLng(delivery.CollectionAddress)
	if err != false {
		delivery.CollectionLat = lat
		delivery.CollectionLon = lng
	} else {
		return u.Message(false, "Failed to create delivery, connection error.")
	}

	lat, lng, err = setLatLng(delivery.DeliveryAddress)
	if err != false {
		delivery.DeliveryLat = lat
		delivery.DeliveryLon = lng
	} else {
		return u.Message(false, "Failed to create delivery, connection error.")
	}

	GetDB().Table("delivery").Create(delivery)

	if delivery.DeliveryID <= 0 {
		return u.Message(false, "Failed to create delivery, connection error.")
	}
	response := u.Message(true, "Delivery created successfully")
	response["delivery"] = delivery
	return response
}

// GetDeliveries gets a list of deliveries from database created by a user with a profileID
func GetDeliveries(profileID uint) []*Delivery {
	deliveries := make([]*Delivery, 0)

	//checks if its a userprofile
	err := GetDB().Table("delivery").Where("userid = ?", profileID).Find(&deliveries).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, delivery := range deliveries {

		fmt.Println(delivery.ExpiryDate)
		delivery.ExpiryDate = formatDateTime(delivery.ExpiryDate)
		fmt.Println(delivery.ExpiryDate)
	}

	return deliveries
}

func formatDateTime(dateTime string) string {

	//removes TZ characters from the string
	new := strings.Replace(dateTime, "T", " ", -1)
	new = strings.Replace(new, "Z", "", -1)

	return new

}

func setLatLng(address string) (float32, float32, bool) {

	geocoder := opencagedata.NewGeocoder(string([]byte(os.Getenv("open_api_key"))))
	result, err := geocoder.Geocode(address, nil)
	if err != nil {
		fmt.Println(err)
		return 0, 0, false
	}
	lat := result.Results[0].Geometry.Latitude
	lng := result.Results[0].Geometry.Longitude
	return lat, lng, true
}

/*
	Radius bound points
	    +
        3
        |
+1 ---  o  --- 2  lng
	    |
		4
		lat
*/

// GetDeliveriesByRadius gets deliveries from the database where deliveries lat lng coordinates are within the
// four radius bounds
func GetDeliveriesByRadius(journey *Journey, radius float32) []*Delivery {
	deliveries := make([]*Delivery, 0)

	lat, lng, error := setLatLng(journey.JourneyStart)
	if error != true {
		return nil
	}
	loc := &Location{Lat: lat, Lng: lng}

	//radius bounds for collection
	cThree := loc.Lat + u.MilesToLat(radius)
	cFour := loc.Lat - u.MilesToLat(radius)
	cOne := loc.Lng + u.MilesToLng(loc.Lat, radius)
	cTwo := loc.Lng - u.MilesToLng(loc.Lat, radius)

	err := GetDB().Table("delivery").Where("accepted = ?", false).Find(&deliveries).Error
	if err != nil {
		return nil
	}

	var deliveriesFound []*Delivery

	for _, delivery := range deliveries {

		fmt.Printf("lat = %f lng = %f\n", delivery.CollectionLat, delivery.CollectionLon)
		fmt.Printf(" lng between : one = %f origin = %f  two = %f\n", cOne, loc.Lng, cTwo)
		fmt.Printf(" lat between : three = %f origin = %f  four = %f\n", cThree, loc.Lat, cFour)

		if delivery.CollectionLon >= cTwo && delivery.CollectionLon <= cOne {
			if delivery.CollectionLat >= cFour && delivery.CollectionLat <= cThree {
				fmt.Println("found")
				deliveriesFound = append(deliveriesFound, delivery)
			}
		}
	}

	for _, delivery := range deliveries {

		//for debuging
		// fmt.Printf("lat = %f lng = %f\n", delivery.CollectionLat, delivery.CollectionLon)
		// fmt.Printf(" lng between : one = %f origin = %f  two = %f\n", cOne, loc.Lng, cTwo)
		// fmt.Printf(" lat between : three = %f origin = %f  four = %f\n", cThree, loc.Lat, cFour)

		if delivery.DeliveryLon >= cTwo && delivery.DeliveryLon <= cOne {
			if delivery.DeliveryLat >= cFour && delivery.DeliveryLat <= cThree {
				fmt.Println("found")
				deliveriesFound = append(deliveriesFound, delivery)
			}
		}
	}

	if len(deliveriesFound) == 0 {
		return nil
	}
	return deliveriesFound
}

// GetDelivery gets a delivery from the database with the deliveryID
func GetDelivery(id uint) *Delivery {
	delivery := &Delivery{}
	GetDB().Table("delivery").Where("deliveryID = ?", id).First(delivery)
	if delivery.Status == "" {
		return nil
	}

	return delivery
}

// AcceptDelivery accepts the delivery by updating the delivery's accepted field to true
//and adds the accepted delivery data to the database
func (acceptedDelivery *Accepted) AcceptDelivery() map[string]interface{} {

	userProfile := GetUserAccountProfile(acceptedDelivery.UserProfileID)
	courierProfile := GetAccountProfile(acceptedDelivery.CourierAccountID)
	delivery := GetDelivery(acceptedDelivery.DeliveryID)

	if userProfile == nil {
		return u.Message(false, "User not found, connection error.")
	}

	if courierProfile == nil {
		return u.Message(false, "Courier not found, connection error.")
	}
	//changes account id to profile id
	acceptedDelivery.CourierAccountID = courierProfile.ProfileID
	tx := GetDB().Begin()

	err := tx.Table("accepted_delivery").Create(acceptedDelivery).Error

	fmt.Println(acceptedDelivery)

	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return u.Message(false, "Failed to accept delivery, connection error.")
	}

	err = tx.Table("delivery").Model(&delivery).Where("deliveryid = ?", delivery.DeliveryID).Update("accepted", true).Error

	if err != nil {
		tx.Rollback()
		return u.Message(false, "Failed to accept delivery, connection error.")
	}

	//send emails

	tx.Commit()
	response := u.Message(true, "Delivery accepted successfully")
	response["delivery"] = acceptedDelivery
	return response
}

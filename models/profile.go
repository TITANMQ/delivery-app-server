package models

import u "backend/utils"

//Profile struct to hold profile data
type Profile struct {
	ProfileID   uint   `gorm:"column:profileid; auto_increment"`
	FirstName   string `json:"firstName" gorm:"column:firstname"`
	LastName    string `json:"lastName" gorm:"column:lastname"`
	AccountID   uint   `json:"accountID" gorm:"column:accountid"`
	AccountType string `json:"accountType" gorm:"-"`
}

//AccountProfile struct holds profile and account data used for get requests
type AccountProfile struct {
	ProfileID   uint   `json:"profileID"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	AccountID   uint   `json:"accountID"`
	AccountType string `json:"accountType"`
}

// Create function creates a profile and saves it to the database
func (profile *Profile) Create() map[string]interface{} {

	if profile.AccountType == "user" {
		GetDB().Table("userprofile").Create(profile)
	} else if profile.AccountType == "courier" {
		GetDB().Table("courierprofile").Create(profile)
	} else {
		return u.Message(false, "Invalid account type.")
	}

	if profile.ProfileID <= 0 {
		return u.Message(false, "Failed to create profile, connection error.")
	}

	response := u.Message(true, "Profile created successfully")
	response["profile"] = profile

	return response

}

// GetProfile gets a profile from database with accountID
func GetProfile(accountID uint) *Profile {
	profile := &Profile{}
	//checks if its a userprofile
	GetDB().Table("userprofile").Where("accountID = ?", accountID).First(profile)
	if profile.ProfileID == 0 { //userprofile not found
		//checks if its a courierprofile
		GetDB().Table("courierprofile").Where("accountID = ?", accountID).First(profile)
		if profile.ProfileID == 0 {
			return nil //user not found
		}
		profile.AccountType = "courier"

	} else {
		profile.AccountType = "user"
	}

	return profile
}

// GetCourierProfile gets a courier profile from database with profileID
func GetCourierProfile(profileID uint) *Profile {
	profile := &Profile{}
	//checks if its a userprofile
	err := GetDB().Table("courierprofile").Where("profileid= ?", profileID).First(profile).Error
	if err != nil { //courierprofile not found
		return nil
	}
	return profile
}

// GetUserProfile gets a user profile from database with profileID
func GetUserProfile(profileID uint) *Profile {
	profile := &Profile{}
	//checks if its a userprofile
	err := GetDB().Table("userprofile").Where("profileID = ?", profileID).First(profile).Error
	if err != nil { //userprofile not found
		return nil
	}
	return profile
}

// GetAccountProfile gets a profile from database with accountID
func GetAccountProfile(accountID uint) *AccountProfile {
	//combines the data from profile and account
	accountProfile := &AccountProfile{}

	acc := GetAccount(uint(accountID))
	if acc == nil {
		return nil
	}

	profile := GetProfile(uint(accountID))
	if profile == nil {
		return nil
	}

	accountProfile.ProfileID = profile.ProfileID
	accountProfile.FirstName = profile.FirstName
	accountProfile.LastName = profile.LastName
	accountProfile.AccountID = profile.AccountID
	accountProfile.AccountType = profile.AccountType
	accountProfile.Email = acc.Email

	return accountProfile
}

// GetUserAccountProfile gets a user profile and account data from database with profileID
func GetUserAccountProfile(profileID uint) *AccountProfile {
	//combines the data from profile and account
	accountProfile := &AccountProfile{}

	profile := GetUserProfile(profileID)
	if profile == nil {
		return nil
	}

	acc := GetAccount(profile.AccountID)
	if acc == nil {
		return nil
	}

	accountProfile.ProfileID = profile.ProfileID
	accountProfile.FirstName = profile.FirstName
	accountProfile.LastName = profile.LastName
	accountProfile.AccountID = profile.AccountID
	accountProfile.AccountType = profile.AccountType
	accountProfile.Email = acc.Email

	return accountProfile
}

// GetCourierAccountProfile gets a courier profile and account data from database with profileID
func GetCourierAccountProfile(profileID uint) *AccountProfile {
	//combines the data from profile and account
	accountProfile := &AccountProfile{}

	profile := GetCourierProfile(profileID)
	if profile == nil {
		return nil
	}

	acc := GetAccount(profile.AccountID)
	if acc == nil {
		return nil
	}

	accountProfile.ProfileID = profile.ProfileID
	accountProfile.FirstName = profile.FirstName
	accountProfile.LastName = profile.LastName
	accountProfile.AccountID = profile.AccountID
	accountProfile.AccountType = profile.AccountType
	accountProfile.Email = acc.Email

	return accountProfile
}

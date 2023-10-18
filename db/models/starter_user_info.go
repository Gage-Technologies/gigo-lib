package models

type UserStart struct {
	Usage             string `json:"usage"`
	Proficiency       string `json:"proficiency"`
	Tags              string `json:"tags"`
	PreferredLanguage string `json:"preferred_language"`
}

var DefaultUserStart = UserStart{
	Usage:             "",
	Proficiency:       "",
	Tags:              "",
	PreferredLanguage: "",
}

var DefaultStartUserInfo = DefaultUserStart

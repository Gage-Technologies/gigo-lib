package models

type AvatarSettings struct {
	TopType         string `json:"topType"`
	AccessoriesType string `json:"accessoriesType"`
	HairColor       string `json:"hairColor"`
	FacialHairType  string `json:"facialHairType"`
	ClotheType      string `json:"clotheType"`
	ClotheColor     string `json:"clotheColor"`
	EyeType         string `json:"eyeType"`
	EyebrowType     string `json:"eyebrowType"`
	MouthType       string `json:"mouthType"`
	AvatarStyle     string `json:"avatarStyle"`
	SkinColor       string `json:"skinColor"`
}

var DefaultAvatarSettings = AvatarSettings{
	TopType:         "ShortHairDreads02",
	AccessoriesType: "Prescription02",
	HairColor:       "BrownDark",
	FacialHairType:  "Blank",
	ClotheType:      "Hoodie",
	ClotheColor:     "PastelBlue",
	EyeType:         "Happy",
	EyebrowType:     "Default",
	MouthType:       "Smile",
	AvatarStyle:     "Circle",
	SkinColor:       "Light",
}

var DefaultAvatarSettingsInfo = DefaultAvatarSettings

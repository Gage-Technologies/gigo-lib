package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type JourneyInfo struct {
	ID               int64               `json:"_id" sql:"_id"`
	UserID           int64               `json:"user_id" sql:"user_id"`
	LearningGoal     string              `json:"learning_goal" sql:"learning_goal"`
	SelectedLanguage ProgrammingLanguage `json:"selected_language" sql:"selected_language"`
	EndGoal          string              `json:"end_goal" sql:"end_goal"`
	ExperienceLevel  string              `json:"experience_level" sql:"experience_level"`
	FamiliarityIDE   string              `json:"familiarity_ide" sql:"familiarity_ide"`
	FamiliarityLinux string              `json:"familiarity_linux" sql:"familiarity_linux"`
	Tried            string              `json:"tried" sql:"tried"`
	TriedOnline      string              `json:"tried_online" sql:"tried_online"`
	AptitudeLevel    string              `json:"aptitude_level" sql:"aptitude_level"`
}

type JourneyInfoSQL struct {
	ID               int64               `json:"_id" sql:"_id"`
	UserID           int64               `json:"user_id" sql:"user_id"`
	LearningGoal     string              `json:"learning_goal" sql:"learning_goal"`
	SelectedLanguage ProgrammingLanguage `json:"selected_language" sql:"selected_language"`
	EndGoal          string              `json:"end_goal" sql:"end_goal"`
	ExperienceLevel  string              `json:"experience_level" sql:"experience_level"`
	FamiliarityIDE   string              `json:"familiarity_ide" sql:"familiarity_ide"`
	FamiliarityLinux string              `json:"familiarity_linux" sql:"familiarity_linux"`
	Tried            string              `json:"tried" sql:"tried"`
	TriedOnline      string              `json:"tried_online" sql:"tried_online"`
	AptitudeLevel    string              `json:"aptitude_level" sql:"aptitude_level"`
}

type JourneyInfoFrontend struct {
	ID               string              `json:"_id" sql:"_id"`
	UserID           string              `json:"user_id" sql:"user_id"`
	LearningGoal     string              `json:"learning_goal" sql:"learning_goal"`
	SelectedLanguage ProgrammingLanguage `json:"selected_language" sql:"selected_language"`
	EndGoal          string              `json:"end_goal" sql:"end_goal"`
	ExperienceLevel  string              `json:"experience_level" sql:"experience_level"`
	FamiliarityIDE   string              `json:"familiarity_ide" sql:"familiarity_ide"`
	FamiliarityLinux string              `json:"familiarity_linux" sql:"familiarity_linux"`
	Tried            string              `json:"tried" sql:"tried"`
	TriedOnline      string              `json:"tried_online" sql:"tried_online"`
	AptitudeLevel    string              `json:"aptitude_level" sql:"aptitude_level"`
}

func CreateJourneyInfo(id int64, userId int64, learningGoal string, language ProgrammingLanguage, endGoal string,
	experienceLevel string, familiarityIDE string, familiarityLinux string, tried string, triedOnline string,
	aptitudeLevel string) (*JourneyInfo, error) {

	return &JourneyInfo{
		ID:               id,
		UserID:           userId,
		LearningGoal:     learningGoal,
		SelectedLanguage: language,
		EndGoal:          endGoal,
		ExperienceLevel:  experienceLevel,
		FamiliarityIDE:   familiarityIDE,
		FamiliarityLinux: familiarityLinux,
		Tried:            tried,
		TriedOnline:      triedOnline,
		AptitudeLevel:    aptitudeLevel,
	}, nil
}

func JourneyInfoFromSQLNative(rows *sql.Rows) (*JourneyInfo, error) {
	var journeyInfoSQL JourneyInfoSQL

	err := sqlstruct.Scan(&journeyInfoSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning journey info in first scan: %v", err))
	}

	return &JourneyInfo{
		ID:               journeyInfoSQL.ID,
		UserID:           journeyInfoSQL.UserID,
		LearningGoal:     journeyInfoSQL.LearningGoal,
		SelectedLanguage: journeyInfoSQL.SelectedLanguage,
		EndGoal:          journeyInfoSQL.EndGoal,
		ExperienceLevel:  journeyInfoSQL.ExperienceLevel,
		FamiliarityIDE:   journeyInfoSQL.FamiliarityIDE,
		FamiliarityLinux: journeyInfoSQL.FamiliarityLinux,
		Tried:            journeyInfoSQL.Tried,
		TriedOnline:      journeyInfoSQL.TriedOnline,
		AptitudeLevel:    journeyInfoSQL.AptitudeLevel,
	}, nil

}

func (i *JourneyInfo) ToFrontend() *JourneyInfoFrontend {
	return &JourneyInfoFrontend{
		ID:               fmt.Sprintf("%d", i.ID),
		UserID:           fmt.Sprintf("%d", i.UserID),
		LearningGoal:     i.LearningGoal,
		SelectedLanguage: i.SelectedLanguage,
		EndGoal:          i.EndGoal,
		ExperienceLevel:  i.ExperienceLevel,
		FamiliarityIDE:   i.FamiliarityIDE,
		FamiliarityLinux: i.FamiliarityLinux,
		Tried:            i.Tried,
		TriedOnline:      i.TriedOnline,
		AptitudeLevel:    i.AptitudeLevel,
	}
}

func (i *JourneyInfo) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into journey_info (_id, user_id, learning_goal, selected_language, end_goal, experience_level, familiarity_ide, familiarity_linux, tried, tried_online, aptitude_level) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.LearningGoal, i.SelectedLanguage, i.EndGoal, i.ExperienceLevel, i.FamiliarityIDE, i.FamiliarityLinux, i.Tried, i.TriedOnline, i.AptitudeLevel},
	})

	return sqlStatements
}

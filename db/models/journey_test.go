package models

import (
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateJourneyInfo(t *testing.T) {
	journey, err := CreateJourneyInfo(1, 69420, "Hobby", "Python", "FullStackDevelopment", "Intermediate",
		"JetBrains", "Intermediate", "Tried", "Tried", "5")
	if err != nil {
		t.Error("\nCreate Journey Info Failed")
		return
	}

	if journey == nil {
		t.Error("\nCreate Journey Info Failed\n    Error: creation returned nil")
		return
	}

	if journey.ID != 1 {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong id")
		return
	}

	if journey.UserID != 69420 {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong user id")
		return
	}

	if journey.LearningGoal != "Hobby" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong learning goal")
		return
	}

	if journey.SelectedLanguage != "Python" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong language")
		return
	}

	if journey.EndGoal != "FullStackDevelopment" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong end goal")
		return
	}

	if journey.ExperienceLevel != "Intermediate" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong experience level")
		return
	}

	if journey.FamiliarityIDE != "JetBrains" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong familiarity ide")
		return
	}

	if journey.FamiliarityLinux != "Intermediate" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong familiarity linux")
		return
	}

	if journey.Tried != "Tried" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong tried")
		return
	}

	if journey.TriedOnline != "Tried" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong tried online")
		return
	}

	if journey.AptitudeLevel != "5" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong aptitude level")
		return
	}

	t.Log("\nCreate Journey Info Succeeded")

}

func TestJourneyInfo_ToSQLNative(t *testing.T) {
	journey, err := CreateJourneyInfo(1, 69420, "Hobby", "Python", "FullStackDevelopment", "Intermediate",
		"JetBrains", "Intermediate", "Tried", "Tried", "5")
	if err != nil {
		t.Error("\nCreate Journey Info Failed")
		return
	}

	statement := journey.ToSQLNative()

	if statement[0].Statement != "insert ignore into journey_info (_id, user_id, learning_goal, selected_language, end_goal, experience_level, familiarity_ide, familiarity_linux, tried, tried_online, aptitude_level) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong statement")
		return
	}

	if len(statement[0].Values) != 11 {
		t.Error("\nCreate Journey Info Failed\n    Error: wrong number of values")
	}

	t.Log("\nCreate Journey Info Succeeded")
}

func TestJourneyInfoFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize Database Failed\n    Error: ", err)
	}

	defer db.DB.Exec("delete from journey_info")

	journey, err := CreateJourneyInfo(1, 69420, "Hobby", "Python", "FullStackDevelopment", "Intermediate",
		"JetBrains", "Intermediate", "Tried", "Tried", "5")
	if err != nil {
		t.Error("\nCreate Journey Info Failed")
		return
	}

	statements := journey.ToSQLNative()

	for _, statement := range statements {
		stmt, err := db.DB.Prepare(statement.Statement)
		if err != nil {
			t.Error("\nPrepare Statement Failed\n    Error: ", err)
		}

		_, err = stmt.Exec(statement.Values...)
		if err != nil {
			t.Error("\nExecute Statement Failed\n    Error: ", err)
		}
	}

	rows, err := db.DB.Query("select * from journey_info where _id = ?", journey.ID)
	if err != nil {
		t.Error("\nCreate Journey Info Failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nCreate Journey Info Failed\n    Error: no rows found")
		return
	}

	j, err := JourneyInfoFromSQLNative(rows)
	if err != nil {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: ", err)
		return
	}

	if j == nil {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: creation returned nil")
		return
	}

	if j.ID != 1 {
		t.Errorf("\nJourney Info From SQL Native Failed\n    Error: wrong id, got: %d", j.ID)
		return
	}

	if j.UserID != 69420 {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong user id")
		return
	}

	if j.LearningGoal != "Hobby" {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong learning goal")
		return
	}

	if j.SelectedLanguage != "Python" {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong language")
		return
	}

	if j.EndGoal != "FullStackDevelopment" {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong end goal")
		return
	}

	if j.ExperienceLevel != "Intermediate" {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong experience level")
		return
	}

	if j.FamiliarityIDE != "JetBrains" {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong familiarity ide")
		return
	}

	if j.FamiliarityLinux != "Intermediate" {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong familiarity linux")
		return
	}

	if j.Tried != "Tried" {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong tried")
		return
	}

	if j.TriedOnline != "Tried" {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong tried online")
		return
	}

	if j.AptitudeLevel != "5" {
		t.Error("\nJourney Info From SQL Native Failed\n    Error: wrong aptitude level")
		return
	}

	t.Log("\nJourney Info From SQL Native Succeeded")

}

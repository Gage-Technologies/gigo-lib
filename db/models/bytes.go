package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kisielk/sqlstruct"
)

type Bytes struct {
	ID                   int64               `json:"_id" sql:"_id"`
	Name                 string              `json:"name" sql:"name"`
	DescriptionEasy      string              `json:"description_easy" sql:"description_easy"`
	DescriptionMedium    string              `json:"description_medium" sql:"description_medium"`
	DescriptionHard      string              `json:"description_hard" sql:"description_hard"`
	OutlineContentEasy   string              `json:"outline_content_easy" sql:"outline_content_easy"`
	OutlineContentMedium string              `json:"outline_content_medium" sql:"outline_content_medium"`
	OutlineContentHard   string              `json:"outline_content_hard" sql:"outline_content_hard"`
	DevStepsEasy         string              `json:"dev_steps_easy" sql:"dev_steps_easy"`
	DevStepsMedium       string              `json:"dev_steps_medium" sql:"dev_steps_medium"`
	DevStepsHard         string              `json:"dev_steps_hard" sql:"dev_steps_hard"`
	QuestionsEasy        []string            `json:"questions_easy" sql:"questions_easy"`
	QuestionsMedium      []string            `json:"questions_medium" sql:"questions_medium"`
	QuestionsHard        []string            `json:"questions_hard" sql:"questions_hard"`
	Lang                 ProgrammingLanguage `json:"lang" sql:"lang"`
	Published            bool                `json:"published" sql:"published"`
	Color                string              `json:"color" sql:"color"`
}

type BytesSQL struct {
	ID                   int64               `json:"_id" sql:"_id"`
	Name                 string              `json:"name" sql:"name"`
	DescriptionEasy      string              `json:"description_easy" sql:"description_easy"`
	DescriptionMedium    string              `json:"description_medium" sql:"description_medium"`
	DescriptionHard      string              `json:"description_hard" sql:"description_hard"`
	OutlineContentEasy   string              `json:"outline_content_easy" sql:"outline_content_easy"`
	OutlineContentMedium string              `json:"outline_content_medium" sql:"outline_content_medium"`
	OutlineContentHard   string              `json:"outline_content_hard" sql:"outline_content_hard"`
	DevStepsEasy         string              `json:"dev_steps_easy" sql:"dev_steps_easy"`
	DevStepsMedium       string              `json:"dev_steps_medium" sql:"dev_steps_medium"`
	DevStepsHard         string              `json:"dev_steps_hard" sql:"dev_steps_hard"`
	QuestionsEasy        []byte              `json:"questions_easy" sql:"questions_easy"`
	QuestionsMedium      []byte              `json:"questions_medium" sql:"questions_medium"`
	QuestionsHard        []byte              `json:"questions_hard" sql:"questions_hard"`
	Lang                 ProgrammingLanguage `json:"lang" sql:"lang"`
	Published            bool                `json:"published" sql:"published"`
	Color                string              `json:"color" sql:"color"`
}

type BytesFrontend struct {
	ID                   string              `json:"_id" sql:"_id"`
	Name                 string              `json:"name" sql:"name"`
	DescriptionEasy      string              `json:"description_easy" sql:"description_easy"`
	DescriptionMedium    string              `json:"description_medium" sql:"description_medium"`
	DescriptionHard      string              `json:"description_hard" sql:"description_hard"`
	OutlineContentEasy   string              `json:"outline_content_easy" sql:"outline_content_easy"`
	OutlineContentMedium string              `json:"outline_content_medium" sql:"outline_content_medium"`
	OutlineContentHard   string              `json:"outline_content_hard" sql:"outline_content_hard"`
	DevStepsEasy         string              `json:"dev_steps_easy" sql:"dev_steps_easy"`
	DevStepsMedium       string              `json:"dev_steps_medium" sql:"dev_steps_medium"`
	DevStepsHard         string              `json:"dev_steps_hard" sql:"dev_steps_hard"`
	QuestionsEasy        []string            `json:"questions_easy" sql:"questions_easy"`
	QuestionsMedium      []string            `json:"questions_medium" sql:"questions_medium"`
	QuestionsHard        []string            `json:"questions_hard" sql:"questions_hard"`
	Lang                 ProgrammingLanguage `json:"lang" sql:"lang"`
	Published            bool                `json:"published" sql:"published"`
	Color                string              `json:"color" sql:"color"`
}

type BytesSearch struct {
	ID          int64               `json:"_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Lang        ProgrammingLanguage `json:"lang"`
	Published   bool                `json:"published"`
}

func CreateBytes(id int64, name string, easyDescription string, mediumDescription string, hardDescription string,
	easyOutlineContent string, mediumOutlineContent string, hardOutlineContent string, easyDevSteps string,
	mediumDevSteps string, hardDevSteps string, easyQuestions []string, mediumQuestions []string, hardQuestions []string,
	lang ProgrammingLanguage, color string) (*Bytes, error) {
	return &Bytes{
		ID:                   id,
		Name:                 name,
		DescriptionEasy:      easyDescription,
		DescriptionHard:      hardDescription,
		DescriptionMedium:    mediumDescription,
		OutlineContentEasy:   easyOutlineContent,
		OutlineContentMedium: mediumOutlineContent,
		OutlineContentHard:   hardOutlineContent,
		DevStepsHard:         hardDevSteps,
		DevStepsMedium:       mediumDevSteps,
		DevStepsEasy:         easyDevSteps,
		QuestionsEasy:        easyQuestions,
		QuestionsMedium:      mediumQuestions,
		QuestionsHard:        hardQuestions,
		Lang:                 lang,
		Color:                color,
	}, nil
}

func BytesFromSQLNative(rows *sql.Rows) (*Bytes, error) {
	bytesSQL := new(BytesSQL)
	err := sqlstruct.Scan(bytesSQL, rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning Bytes info in first scan: %v", err)
	}

	var questionsEasy []string
	if len(bytesSQL.QuestionsEasy) > 0 {
		err = json.Unmarshal(bytesSQL.QuestionsEasy, &questionsEasy)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling Questions Easy JSON into slice of strings: %v", err)
		}
	} else {
		questionsEasy = make([]string, 0)
	}

	var questionsMedium []string
	if len(bytesSQL.QuestionsMedium) > 0 {
		err = json.Unmarshal(bytesSQL.QuestionsMedium, &questionsMedium)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling Questions Medium JSON into slice of strings: %v", err)
		}
	} else {
		questionsMedium = make([]string, 0)
	}

	var questionsHard []string
	if len(bytesSQL.QuestionsHard) > 0 {
		err = json.Unmarshal(bytesSQL.QuestionsHard, &questionsHard)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling Questions Hard JSON into slice of strings: %v", err)
		}
	} else {
		questionsHard = make([]string, 0)
	}

	return &Bytes{
		ID:                   bytesSQL.ID,
		Name:                 bytesSQL.Name,
		DescriptionEasy:      bytesSQL.DescriptionEasy,
		DescriptionMedium:    bytesSQL.DescriptionMedium,
		DescriptionHard:      bytesSQL.DescriptionHard,
		OutlineContentEasy:   bytesSQL.OutlineContentEasy,
		OutlineContentMedium: bytesSQL.OutlineContentMedium,
		OutlineContentHard:   bytesSQL.OutlineContentHard,
		DevStepsEasy:         bytesSQL.DevStepsEasy,
		DevStepsMedium:       bytesSQL.DevStepsMedium,
		DevStepsHard:         bytesSQL.DevStepsHard,
		Lang:                 bytesSQL.Lang,
		QuestionsEasy:        questionsEasy,
		QuestionsMedium:      questionsMedium,
		QuestionsHard:        questionsHard,
		Published:            bytesSQL.Published,
		Color:                bytesSQL.Color,
	}, nil
}

func (b *Bytes) ToFrontend() *BytesFrontend {
	return &BytesFrontend{
		ID:                   fmt.Sprintf("%d", b.ID),
		Name:                 b.Name,
		DescriptionEasy:      b.DescriptionEasy,
		DescriptionMedium:    b.DescriptionMedium,
		DescriptionHard:      b.DescriptionHard,
		OutlineContentEasy:   b.OutlineContentEasy,
		OutlineContentMedium: b.OutlineContentMedium,
		OutlineContentHard:   b.OutlineContentHard,
		DevStepsEasy:         b.DevStepsEasy,
		DevStepsMedium:       b.DevStepsMedium,
		DevStepsHard:         b.DevStepsHard,
		QuestionsEasy:        b.QuestionsEasy,
		QuestionsMedium:      b.QuestionsMedium,
		QuestionsHard:        b.QuestionsHard,
		Lang:                 b.Lang,
		Published:            b.Published,
		Color:                b.Color,
	}
}

func (b *Bytes) ToSearch() *BytesSearch {
	return &BytesSearch{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.DescriptionMedium,
		Lang:        b.Lang,
		Published:   b.Published,
	}
}

func (b *BytesSearch) ToFrontend() *BytesFrontend {
	return &BytesFrontend{
		ID:                   fmt.Sprintf("%d", b.ID),
		Name:                 b.Name,
		DescriptionMedium:    b.Description,
		Lang:                 b.Lang,
		Published:            b.Published,
	}
}

func (b *Bytes) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)

	var questionsEasyJSON []byte
	if len(b.QuestionsEasy) > 0 {
		var err error
		questionsEasyJSON, err = json.Marshal(b.QuestionsEasy)
		if err != nil {
			return nil, fmt.Errorf("error marshaling Questions Easy JSON: %v", err)
		}
	}

	var questionsMediumJSON []byte
	if len(b.QuestionsMedium) > 0 {
		var err error
		questionsMediumJSON, err = json.Marshal(b.QuestionsMedium)
		if err != nil {
			return nil, fmt.Errorf("error marshaling Questions Medium JSON: %v", err)
		}
	}

	var questionsHardJSON []byte
	if len(b.QuestionsHard) > 0 {
		var err error
		questionsHardJSON, err = json.Marshal(b.QuestionsHard)
		if err != nil {
			return nil, fmt.Errorf("error marshaling Questions Hard JSON: %v", err)
		}
	}

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into bytes(_id, name, description_easy, description_medium, description_hard, outline_content_easy, outline_content_medium, outline_content_hard, dev_steps_easy, dev_steps_medium, dev_steps_hard, lang, published, color, questions_easy, questions_medium, questions_hard) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);",
		Values:    []interface{}{b.ID, b.Name, b.DescriptionEasy, b.DescriptionMedium, b.DescriptionHard, b.OutlineContentEasy, b.OutlineContentMedium, b.OutlineContentHard, b.DevStepsEasy, b.DevStepsMedium, b.DevStepsHard, b.Lang, b.Published, b.Color, questionsEasyJSON, questionsMediumJSON, questionsHardJSON},
	})

	return sqlStatements, nil
}

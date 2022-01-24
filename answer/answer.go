package answer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var answersFile = "answers.json"

type Answer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetAnswerValue retrieves an answer's value for a given answer's key
// If the answer does not exist it will exit with an error
func GetAnswerValue(key string) (string, error) {
	value, ok := AllAnswersIndexed[key]
	if !ok {
		return "", fmt.Errorf("getAnswerValue: answer '%s' does not exist", key)
	}
	return value, nil
}

// CreateAnswer creates an answer, if the answer already exists it exits with an error
func CreateAnswer(ans Answer) error {
	_, ok := AllAnswersIndexed[ans.Key]
	if ok {
		return fmt.Errorf("createAnswer: answer '%s' already exists", ans.Key)
	}

	AllAnswersIndexed[ans.Key] = ans.Value
	err := AllAnswersIndexed.SaveAnswers() // refresh all answers
	if err != nil {
		return err
	}
	err = AllEventsIndexed.CreateEvent(Create, ans) // create event
	if err != nil {
		return err
	}
	return nil
}

// Deletevalue will try to delete an existing answer given its key
// If the requested answer does not exist it will error out
func DeleteAnswer(key string) error {
	value, ok := AllAnswersIndexed[key]
	if !ok {
		return fmt.Errorf("deleteAnswer: answer '%s' does not exist", key)
	}
	delete(AllAnswersIndexed, key)
	err := AllAnswersIndexed.SaveAnswers()
	if err != nil {
		return err
	} // refresh all answers
	err = AllEventsIndexed.CreateEvent(Delete, Answer{Key: key, Value: value}) // delete event
	if err != nil {
		return err
	}
	return nil
}

// EditAnswer will try to update an existing answer given its key
// If the requested answer does not exist it will error out
func EditAnswer(ans Answer) error {
	_, ok := AllAnswersIndexed[ans.Key]
	if !ok {
		return fmt.Errorf("editAnswer: answer '%s' does not exist", ans.Key)
	}
	AllAnswersIndexed[ans.Key] = ans.Value
	err := AllAnswersIndexed.SaveAnswers() // refresh all answers
	if err != nil {
		return err
	}
	err = AllEventsIndexed.CreateEvent(Update, ans) // update event
	if err != nil {
		return err
	}
	return nil
}

type MapOfAnswers map[string]string

var AllAnswersIndexed MapOfAnswers

// SaveAnswers will store all existing answers into a json file for preservation
func (mapAnswers MapOfAnswers) SaveAnswers() error {
	allAnswers := make([]Answer, 0)
	for key, value := range mapAnswers {
		allAnswers = append(allAnswers, Answer{Key: key, Value: value})
	}

	file, err := json.MarshalIndent(allAnswers, "", "")
	if err != nil {
		return fmt.Errorf("saveAnswers: %v", err)
	}
	err = ioutil.WriteFile(answersFile, file, 0644)
	if err != nil {
		return fmt.Errorf("saveAnswers: %v", err)
	}
	return nil
}

// InitAnswers will populate internal answers map using the json answers file
// If file is not found or it is not valid, an empty map will be created instead
func InitAnswers() {
	AllAnswersIndexed = make(MapOfAnswers)

	body, err := ioutil.ReadFile(answersFile)
	if err != nil {
		return
	}

	var allAnswers []Answer
	err = json.Unmarshal(body, &allAnswers)
	if err != nil {
		return
	}

	for _, ans := range allAnswers {
		AllAnswersIndexed[ans.Key] = ans.Value
	}
}

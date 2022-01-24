package answer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var answersFile = "answers.json"

type Answer struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// GetValue retrieves an answer's value for a given answer's key
// If the answer does not exist it will exit with an error
func GetValue(key string) (string, error) {
	value, ok := AllAnswersIndexed[key]
	if !ok {
		return "", fmt.Errorf("getValue: answer '%s' does not exist", key)
	}
	return value, nil
}

// PostValue creates an answer, if the answer already exists it exits with an error
func PostValue(ans Answer) error {
	_, ok := AllAnswersIndexed[ans.Key]
	if ok {
		return fmt.Errorf("postValue: answer '%s' already exists", ans.Key)
	}

	AllAnswersIndexed[ans.Key] = ans.Value
	AllAnswersIndexed.SaveAnswers() // refresh all answers
	return nil
}

// Deletevalue will try to delete an existing answer given its key
// If the requested answer does not exist it will error out
func DeleteValue(key string) error {
	_, ok := AllAnswersIndexed[key]
	if !ok {
		return fmt.Errorf("deleteValue: answer '%s' does not exist", key)
	}
	delete(AllAnswersIndexed, key)
	return nil
}

// EditValue will try to update an existing answer given its key
// If the requested answer does not exist it will error out
func EditValue(ans Answer) error {
	_, ok := AllAnswersIndexed[ans.Key]
	if !ok {
		return fmt.Errorf("editValue: answer '%s' does not exist", ans.Key)
	}
	AllAnswersIndexed[ans.Key] = ans.Value
	AllAnswersIndexed.SaveAnswers()
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
func InitAnswers() error {
	body, _ := ioutil.ReadFile(answersFile)

	var allAnswers []Answer
	err := json.Unmarshal(body, &allAnswers)
	if err != nil {
		return fmt.Errorf("initAnswers: %v", err)
	}

	AllAnswersIndexed = make(MapOfAnswers)

	for _, ans := range allAnswers {
		AllAnswersIndexed[ans.Key] = ans.Value
	}

	return nil
}

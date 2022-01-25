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

type MapOfAnswers map[string]string

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

// InitAnswers will populate an answers map using the json answers file
// If the file is not found or it is not valid, an empty map will be created instead
func InitAnswers() MapOfAnswers {
	allAnswersIndexed := make(MapOfAnswers)

	body, err := ioutil.ReadFile(answersFile)
	if err != nil {
		return allAnswersIndexed
	}

	var allAnswers []Answer
	err = json.Unmarshal(body, &allAnswers)
	if err != nil {
		return allAnswersIndexed
	}

	for _, ans := range allAnswers {
		allAnswersIndexed[ans.Key] = ans.Value
	}

	return allAnswersIndexed
}

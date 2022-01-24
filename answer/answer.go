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

func GetValue(key string) string {
	return AllAnswersIndexed[key]
}

func PostValue(ans Answer) {
	AllAnswersIndexed[ans.Key] = ans.Value
	AllAnswersIndexed.SaveAnswers()
}

func DeleteValue(key string) {
	_, ok := AllAnswersIndexed[key]
	if ok {
		delete(AllAnswersIndexed, key)
	}
}

func EditValue(ans Answer) error {
	_, ok := AllAnswersIndexed[ans.Key]
	if !ok {
		return fmt.Errorf("could not edit value for non existing answer")
	}
	AllAnswersIndexed[ans.Key] = ans.Value
	AllAnswersIndexed.SaveAnswers()
	return nil
}

type MapOfAnswers map[string]string

var AllAnswersIndexed MapOfAnswers

func (mapAnswers MapOfAnswers) SaveAnswers() {
	allAnswers := make([]Answer, 0)
	for key, value := range mapAnswers {
		allAnswers = append(allAnswers, Answer{Key: key, Value: value})
	}

	file, _ := json.MarshalIndent(allAnswers, "", "") // TODO err handling
	ioutil.WriteFile(answersFile, file, 0644)
}

func InitAnswers() { // TODO err handling
	body, _ := ioutil.ReadFile(answersFile)
	var allAnswers []Answer
	json.Unmarshal(body, &allAnswers)

	AllAnswersIndexed = make(MapOfAnswers)

	for _, ans := range allAnswers {
		AllAnswersIndexed[ans.Key] = ans.Value
	}

	fmt.Println(AllAnswersIndexed)
}

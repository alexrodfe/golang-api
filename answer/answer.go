package answer

import "fmt"

type Answer struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func GetValue(key string) string {
	return AllAnswersIndexed[key]
}

func PostValue(ans Answer) {
	AllAnswersIndexed[ans.Key] = ans.Value
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
	return nil
}

type MapOfAnswers map[string]string

var AllAnswersIndexed MapOfAnswers

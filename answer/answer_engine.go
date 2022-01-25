package answer

import "fmt"

type AnswerEngine struct {
	answersMap *MapOfAnswers
	eventsMap  *MapOfEvents
}

func NewAnswerEngine(answersMap *MapOfAnswers, eventsMap *MapOfEvents) AnswerEngine {
	return AnswerEngine{
		answersMap: answersMap,
		eventsMap:  eventsMap,
	}
}

// GetAnswerValue retrieves an answer's value for a given answer's key
// If the answer does not exist it will exit with an error
func (anse AnswerEngine) GetAnswerValue(key string) (string, error) {
	value, ok := (*anse.answersMap)[key]
	if !ok {
		return "", fmt.Errorf("getAnswerValue: answer '%s' does not exist", key)
	}
	return value, nil
}

// CreateAnswer creates an answer, if the answer already exists it exits with an error
func (anse AnswerEngine) CreateAnswer(ans Answer) error {
	_, ok := (*anse.answersMap)[ans.Key]
	if ok {
		return fmt.Errorf("createAnswer: answer '%s' already exists", ans.Key)
	}

	(*anse.answersMap)[ans.Key] = ans.Value
	err := anse.answersMap.SaveAnswers() // refresh all answers
	if err != nil {
		return err
	}
	err = anse.eventsMap.CreateEvent(Create, ans) // create event
	if err != nil {
		return err
	}
	return nil
}

// Deletevalue will try to delete an existing answer given its key
// If the requested answer does not exist it will error out
func (anse AnswerEngine) DeleteAnswer(key string) error {
	value, ok := (*anse.answersMap)[key]
	if !ok {
		return fmt.Errorf("deleteAnswer: answer '%s' does not exist", key)
	}
	delete((*anse.answersMap), key)
	err := anse.answersMap.SaveAnswers()
	if err != nil {
		return err
	} // refresh all answers
	err = anse.eventsMap.CreateEvent(Delete, Answer{Key: key, Value: value}) // delete event
	if err != nil {
		return err
	}
	return nil
}

// EditAnswer will try to update an existing answer given its key
// If the requested answer does not exist it will error out
func (anse AnswerEngine) EditAnswer(ans Answer) error {
	_, ok := (*anse.answersMap)[ans.Key]
	if !ok {
		return fmt.Errorf("editAnswer: answer '%s' does not exist", ans.Key)
	}
	(*anse.answersMap)[ans.Key] = ans.Value
	err := anse.answersMap.SaveAnswers() // refresh all answers
	if err != nil {
		return err
	}
	err = anse.eventsMap.CreateEvent(Update, ans) // update event
	if err != nil {
		return err
	}
	return nil
}

package answer

import (
	"github.com/stretchr/testify/suite"
)

var answersFileTest = "answers_test.json"
var eventsFileTest = "events_test.json"

type AnswerTestSuite struct {
	suite.Suite
	allAnswersIndexedTest MapOfAnswers
	allEventsIndexedTest  MapOfEvents
}

func (suite *AnswerTestSuite) SetupSuite() {

}

func (suite *AnswerTestSuite) SetupTest() {
	suite.allAnswersIndexedTest = make(MapOfAnswers)
	suite.allEventsIndexedTest = make(MapOfEvents)
}

func (suite *AnswerTestSuite) TearDownTest() {
}

func (suite *AnswerTestSuite) TestGetAnswerValue() {

}

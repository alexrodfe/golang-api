package answer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

var AllAnswersIndexedMock MapOfAnswers
var AllEventsIndexedMock MapOfEvents

type AnswerTestSuite struct {
	suite.Suite
	anse AnswerEngine
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(AnswerTestSuite))
}

func (suite *AnswerTestSuite) SetupSuite() {
	AllAnswersIndexedMock = make(MapOfAnswers)
	AllEventsIndexedMock = make(MapOfEvents)

	suite.anse = NewAnswerEngine(&AllAnswersIndexedMock, &AllEventsIndexedMock)
}

func (suite *AnswerTestSuite) SetupTest() {
	(*suite.anse.answersMap) = make(MapOfAnswers)
	(*suite.anse.eventsMap) = make(MapOfEvents)
}

func (suite *AnswerTestSuite) TearDownTest() {
	os.Remove(answersFile)
	os.Remove(eventsFile)
}

func createMockanswer() Answer {
	return Answer{
		Key:   "test",
		Value: "ok",
	}
}

func (suite *AnswerTestSuite) TestGetAnswerValue() {
	ans := createMockanswer()

	(*suite.anse.answersMap)[ans.Key] = ans.Value

	testValue, err := suite.anse.GetAnswerValue(ans.Key)
	suite.Require().NoError(err)
	suite.Assert().Equal(ans.Value, testValue)

	_, err = suite.anse.GetAnswerValue("wrong")
	suite.Require().Error(err)
}

func (suite *AnswerTestSuite) TestGetLatestAnswerValue() {
	ans := createMockanswer()
	editValue := "also ok"

	err := suite.anse.CreateAnswer(ans)
	suite.Require().NoError(err)
	ans.Value = editValue
	err = suite.anse.EditAnswer(ans)
	suite.Require().NoError(err)
	testValue, err := suite.anse.GetAnswerValue(ans.Key)
	suite.Require().NoError(err)

	suite.Assert().Equal(editValue, testValue)
}

func (suite *AnswerTestSuite) TestCreateAnswer() {
	ans := createMockanswer()

	_, err := suite.anse.GetAnswerValue(ans.Key)
	suite.Require().Error(err)

	err = suite.anse.CreateAnswer(ans)
	suite.Require().NoError(err)

	_, err = suite.anse.GetAnswerValue(ans.Key)
	suite.Require().NoError(err)

	err = suite.anse.CreateAnswer(ans)
	suite.Require().Error(err)
}

func (suite *AnswerTestSuite) TestDeleteAnswer() {
	ans := createMockanswer()

	_, err := suite.anse.GetAnswerValue(ans.Key)
	suite.Require().Error(err)

	err = suite.anse.CreateAnswer(ans)
	suite.Require().NoError(err)

	_, err = suite.anse.GetAnswerValue(ans.Key)
	suite.Require().NoError(err)

	err = suite.anse.DeleteAnswer(ans.Key)
	suite.Require().NoError(err)

	_, err = suite.anse.GetAnswerValue(ans.Key)
	suite.Require().Error(err)
}

func (suite *AnswerTestSuite) TestEditAnswer() {
	ans := createMockanswer()

	_, err := suite.anse.GetAnswerValue(ans.Key)
	suite.Require().Error(err)

	err = suite.anse.CreateAnswer(ans)
	suite.Require().NoError(err)

	_, err = suite.anse.GetAnswerValue(ans.Key)
	suite.Require().NoError(err)

	editValue := "test"
	ans.Value = editValue
	err = suite.anse.EditAnswer(ans)
	suite.Require().NoError(err)

	testValue, err := suite.anse.GetAnswerValue(ans.Key)
	suite.Require().NoError(err)
	suite.Equal(editValue, testValue)
}

func (suite *AnswerTestSuite) TestGetAnswerHistory() {
	ans := createMockanswer()

	_, err := suite.anse.GetAnswerHistory(ans.Key)
	suite.Require().Error(err)

	err = suite.anse.CreateAnswer(ans)
	suite.Require().NoError(err)

	editValue := "test"
	ans.Value = editValue
	err = suite.anse.EditAnswer(ans)
	suite.Require().NoError(err)

	err = suite.anse.DeleteAnswer(ans.Key)
	suite.Require().NoError(err)

	events, err := suite.anse.GetAnswerHistory(ans.Key)
	suite.Require().NoError(err)
	suite.Require().Len(events, 3)

	// chronological order
	suite.Equal(Create, events[0].Event)
	suite.Equal(Update, events[1].Event)
	suite.Equal(Delete, events[2].Event)
}

func (suite *AnswerTestSuite) TestFlow() {
	ans := createMockanswer()

	// correct flow create → update → delete → create → update
	err := suite.anse.CreateAnswer(ans)
	suite.Require().NoError(err)

	err = suite.anse.EditAnswer(ans)
	suite.Require().NoError(err)

	err = suite.anse.DeleteAnswer(ans.Key)
	suite.Require().NoError(err)

	err = suite.anse.CreateAnswer(ans)
	suite.Require().NoError(err)

	err = suite.anse.EditAnswer(ans)
	suite.Require().NoError(err)

	// incorrect flow create → delete → update
	err = suite.anse.DeleteAnswer(ans.Key)
	suite.Require().NoError(err)

	err = suite.anse.EditAnswer(ans)
	suite.Require().Error(err)
}

func (suite *AnswerTestSuite) TestAnswerSaving() {
	answersMap := InitAnswers()
	suite.Empty(answersMap)

	ans := createMockanswer()
	ans2 := createMockanswer()
	ans2.Key = "test2"

	err := suite.anse.CreateAnswer(ans)
	suite.Require().NoError(err)
	err = suite.anse.CreateAnswer(ans2)
	suite.Require().NoError(err)
	ans2.Value = "also ok"
	err = suite.anse.EditAnswer(ans2)
	suite.Require().NoError(err)

	suite.Len((*suite.anse.answersMap), 2)

	answersMap = InitAnswers()
	suite.Len(answersMap, 2)
}

func (suite *AnswerTestSuite) TestEventSaving() {
	eventsMap := InitEvents()
	suite.Empty(eventsMap)

	ans := createMockanswer()
	ans2 := createMockanswer()
	ans2.Key = "test2"

	err := suite.anse.CreateAnswer(ans)
	suite.Require().NoError(err)
	err = suite.anse.CreateAnswer(ans2)
	suite.Require().NoError(err)
	ans2.Value = "also ok"
	err = suite.anse.EditAnswer(ans2)
	suite.Require().NoError(err)

	suite.Len((*suite.anse.eventsMap), 2)
	suite.Len((*suite.anse.eventsMap)[ans.Key], 1)
	suite.Len((*suite.anse.eventsMap)[ans2.Key], 2)

	eventsMap = InitEvents()
	suite.Len(eventsMap, 2)
	suite.Require().Len(eventsMap[ans.Key], 1)
	suite.Require().Len(eventsMap[ans2.Key], 2)

	suite.Equal(Create, eventsMap[ans.Key][0].Event)
	suite.Equal(Create, eventsMap[ans2.Key][0].Event)
	suite.Equal(Update, eventsMap[ans2.Key][1].Event)
	suite.Equal(ans, eventsMap[ans.Key][0].Data)
	suite.NotEqual(ans2, eventsMap[ans2.Key][0].Data)
	suite.Equal(ans2, eventsMap[ans2.Key][1].Data)
}

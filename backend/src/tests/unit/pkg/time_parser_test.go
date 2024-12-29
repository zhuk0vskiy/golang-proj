package pkg_test

import (

	"backend/src/pkg/time_parser"
	// utils2 "backend/src/tests/utils"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"time"
)

type TimeParserSuite struct {
	suite.Suite
}

func (suite *TimeParserSuite) TestStringToDate01(t provider.T) {
	t.Title("[TimeParser] correct string")
	t.Tags("time_parser", "pass")
	t.Parallel()

	t.WithNewStep("pass", func(sCtx provider.StepCtx) {
		
		string := "2020-09-01 12:45:55";
		expDate := time.Date(2020, 9, 1, 12, 45, 55, 0, time.UTC)

		date, err := time_parser.StringToDate(string)

		sCtx.Assert().Equal(expDate, date)
		sCtx.Assert().NoError(err)
	})
}

func (suite *TimeParserSuite) TestStringToDate02(t provider.T) {
	t.Title("[TimeParser] incorrect string")
	t.Tags("time_parser", "pass")
	t.Parallel()

	t.WithNewStep("pass", func(sCtx provider.StepCtx) {
		
		string := "2020-09-01 12:45";
		// expDate := time.Date(2020, 9, 1, 12, 45, 55, 0, time.UTC)

		_, err := time_parser.StringToDate(string)

		// sCtx.Assert().Equal(expDate, date)
		sCtx.Assert().Error(err)
	})
}

func (suite *TimeParserSuite) TestStringToDate03(t provider.T) {
	t.Title("[TimeParser] incorrect string")
	t.Tags("time_parser", "pass")
	t.Parallel()

	t.WithNewStep("pass", func(sCtx provider.StepCtx) {
		
		string := "";
		// expDate := time.Date(2020, 9, 1, 12, 45, 55, 0, time.UTC)

		_, err := time_parser.StringToDate(string)

		// sCtx.Assert().Equal(expDate, date)
		sCtx.Assert().Error(err)
	})
}
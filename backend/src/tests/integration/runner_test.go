package integration

import (
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"sync"
	"testing"
)

func TestRunner(t *testing.T) {
	t.Parallel()

	wg := &sync.WaitGroup{}
	suits := []runner.TestSuite{
		// &ITCompanySuite{},
		// &ITAuthSuite{},
		// &ITFinReportSuite{},
		// &ITActFieldSuite{},
		// &ITUserSuite{},
	}
	wg.Add(len(suits))

	for _, s := range suits {
		go func(s runner.TestSuite) {
			suite.RunSuite(t, s)
			wg.Done()
		}(s)
	}

	wg.Wait()
}
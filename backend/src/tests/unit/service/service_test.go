package service_test

import (
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"sync"
	"testing"
	"backend/src/pkg/base"
)

func Test_Runner(t *testing.T) {
	//db, ids := utils.NewTestStorage()
	//defer utils.DropTestStorage(db)
	//
	tm, err := base.GenerateAuthToken("1", "2", "3", "4")
	if err != nil {
		panic(err)
	}

	t.Parallel()

	wg := &sync.WaitGroup{}
	suits := []runner.TestSuite{
		&AuthSuite{
			JwtKey: tm,
		},
		&EquipmentSuite{},
		&ProducerSuite{},
		&InstrumentalistSuite{},
		&RoomSuite{},
		//&ReserveSuite{},
		&StudioSuite{},
	}
	wg.Add(len(suits))

	for _, s := range suits {
		go func() {
			suite.RunSuite(t, s)
			wg.Done()
		}()
	}

	wg.Wait()
}

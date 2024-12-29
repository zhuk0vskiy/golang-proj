package service_test

import (
	"backend/src/internal/model"
	"backend/src/internal/repository/interface/mocks"
	"backend/src/internal/service/impl"
	utils2 "backend/src/tests/utils"
	"context"
	// "fmt"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"time"
)

type ReserveSuite struct {
	suite.Suite
}

func (suite *ReserveSuite) TestReserveAdd01(t provider.T) {
	t.Title("[Add] success")
	t.Tags("reserve", "add")
	t.Parallel()

	t.WithNewStep("successfully add", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		request := utils2.ReserveFabric{Id: 1}.CorrectReserveAdd()

		reserveRepo := new(mocks.IReserveRepository)
		roomRepo := new(mocks.IRoomRepository)
		produceRepo := new(mocks.IProducerRepository)
		instrumentalistRepo := new(mocks.IInstrumentalistRepository)

		reserveRepo.On("Add", ctx, request).Return(
			nil,
		).Once()
		roomRepo.On("Add", ctx, request).Return(
			nil,
		).Once()
		produceRepo.On("Add", ctx, request).Return(
			nil,
		).Once()
		instrumentalistRepo.On("Add", ctx, request).Return(
			nil,
		).Once()

		logger := utils2.NewMockLogger()

		service := impl.NewReserveService(logger, reserveRepo, roomRepo, produceRepo, instrumentalistRepo)
		sCtx.WithNewParameters("ctx", ctx, "request", request)

		err := service.Add(ctx, request)

		sCtx.Assert().NoError(err)
	})
}

func (suite *ReserveSuite) TestReserveAdd02(t provider.T) {
	t.Title("[Add] failed")
	t.Tags("reserve", "add")
	t.Parallel()

	t.WithNewStep("failed to add", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		request := utils2.ReserveFabric{Id: 1}.IncorrectReserveAdd()


		reserveRepo := new(mocks.IReserveRepository)
		roomRepo := new(mocks.IRoomRepository)
		produceRepo := new(mocks.IProducerRepository)
		instrumentalistRepo := new(mocks.IInstrumentalistRepository)

		reserveRepo.On("Add", ctx, request).Return(
			nil,
		).Once()
		roomRepo.On("Add", ctx, request).Return(
			nil,
		).Once()
		produceRepo.On("Add", ctx, request).Return(
			nil,
		).Once()
		instrumentalistRepo.On("Add", ctx, request).Return(
			nil,
		).Once()

		logger := utils2.NewMockLogger()

		service := impl.NewReserveService(logger, reserveRepo, roomRepo, produceRepo, instrumentalistRepo)
		sCtx.WithNewParameters("ctx", ctx, "request", request)

		err := service.Add(ctx, request)

		sCtx.Assert().Error(err)
	})
}

func (suite *ReserveSuite) TestReserveGetAll01(t provider.T) {
	t.Title("[GetAll] success")
	t.Tags("reserve", "getAll")
	t.Parallel()

	t.WithNewStep("success to getAll", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		request := utils2.ReserveFabric{Id: 1}.CorrectReserveGetAll()

		reserves := []*model.Reserve{
			{
				UserId:            1,
				RoomId:            1,
				ProducerId:        1,
				InstrumentalistId: 1,
				TimeInterval: &model.TimeInterval{
					StartTime: time.Date(2024, 4, 14, 12, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2024, 4, 14, 13, 0, 0, 0, time.UTC),
				},
			},
		}


		reserveRepo := new(mocks.IReserveRepository)
		roomRepo := new(mocks.IRoomRepository)
		produceRepo := new(mocks.IProducerRepository)
		instrumentalistRepo := new(mocks.IInstrumentalistRepository)

		reserveRepo.On("Add", ctx, request).Return(
			nil,
		).Once()
		roomRepo.On("Add", ctx, request).Return(
			nil,
		).Once()
		produceRepo.On("Add", ctx, request).Return(
			nil,
		).Once()
		instrumentalistRepo.On("Add", ctx, request).Return(
			nil,
		).Once()

		logger := utils2.NewMockLogger()

		service := impl.NewReserveService(logger, reserveRepo, roomRepo, produceRepo, instrumentalistRepo)
		sCtx.WithNewParameters("ctx", ctx, "request", request)

		reserves, err := service.GetAll(request)

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(reserves)

	})
}

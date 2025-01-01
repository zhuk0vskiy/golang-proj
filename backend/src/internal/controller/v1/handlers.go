package v1

import (
	"backend/src/internal/app"
	"backend/src/internal/controller"
	"backend/src/internal/model"
	"backend/src/internal/model/dto"
	"backend/src/pkg/time_parser"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"io"
	"github.com/go-chi/chi/v5"
)

func LoginHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "LogInHandler"
		//start := time.Now()

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//defer func() {
		//	controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
		//}()

		var req LogInRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		ua := &dto.LogInRequest{
			Login:    req.Login,
			Password: req.Password,
		}
		token, err := app.AuthSvc.LogIn(r.Context(), ua)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		cookie := http.Cookie{
			Name:    "access_token",
			Value:   token,
			Path:    "/",
			Secure:  true,
			Expires: time.Now().Add(3600 * 24 * time.Second),
		}

		http.SetCookie(w, &cookie)
		controller.SuccessResponse(wrappedWriter, http.StatusOK, map[string]string{"token": token})
	}
}

func ValidationHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "ValidationHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		var req ValidationRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		timeInterval, err := time_parser.ToTime(req.Date, req.StartHour, req.EndHour)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}
		rooms, equipments, producers, instrumentalists, err := app.ValidateTimeSvc.GetSuitableStuff(
			r.Context(),
			&dto.GetSuitableStuffRequest{
				StudioId:     int64(req.StudioId),
				TimeInterval: timeInterval,
			},
		)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		equipmentsExt := &ValidateEquipment{
			Microphones: equipments[0],
			Guitars:     equipments[1],
		}
		//equipmentsExt := make([]*ValidateEquipment, 0)
		//for equipment := range equipments {
		//	equipmentsExt = append(equipmentsExt, &ValidateEquipment{
		//		Microphones: equipment
		//	})
		//}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, ValidateResponse{
			Rooms:            rooms,
			Producers:        producers,
			Instrumentalists: instrumentalists,
			Equipments:       equipmentsExt,
		})

	}
}

func AddReserveHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "AddReserveHandler"

		var req AddReserveRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		id, err := controller.GetStringClaimFromJWT(r.Context(), "id")
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		timeInterval, err := time_parser.ToTime(req.Date, req.StartHour, req.EndHour)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		err = app.ReserveSvc.Add(r.Context(), &dto.AddReserveRequest{
			//UserId: int64(),
			//UserId:            1,
			UserId:            int64(idInt),
			RoomId:            req.RoomId,
			ProducerId:        req.ProducerId,
			InstrumentalistId: req.InstrumentalistId,
			TimeInterval:      timeInterval,
			EquipmentId:       req.EquipmentsId,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func DeleteReserveHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "DeleteReserveHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.ReserveSvc.Delete(
			&dto.DeleteReserveRequest{
				Id: int64(idInt),
			},
		)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func SignInHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "SignInHandler"

		var req SignInRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.AuthSvc.SignIn(&dto.SignInRequest{
			Login:      req.Login,
			Password:   req.Password,
			FirstName:  req.FirstName,
			SecondName: req.SecondName,
			ThirdName:  req.ThirdName,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func GetUserReservesHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UserReservesHandler"
		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id, err := controller.GetStringClaimFromJWT(r.Context(), "id")
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		reserves, err := app.UserSvc.GetReserves(
			&dto.GetUserReservesRequest{
				Id: int64(idInt),
			},
		)

		var reservesExt []*model.ReserveExt
		//var equipmentsId [][]int64
		//reservesExt := make([]*model.ReserveExt, len(reserves))

		for i := range reserves {
			equipments, err := app.EquipmentSvc.GetByReserve(&dto.GetEquipmentByReserveRequest{
				ReserveId: reserves[i].Id,
			})
			equipmentsId := make([]int64, 0)

			for t := range equipments {
				equipmentsId = append(equipmentsId, equipments[t].Id)
				//	equipmentsId = append(equipmentsId, )
			}
			if err != nil {
				controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
				return
			}

			reservesExt = append(reservesExt, &model.ReserveExt{
				Id:                reserves[i].Id,
				UserId:            reserves[i].UserId,
				RoomId:            reserves[i].RoomId,
				ProducerId:        reserves[i].ProducerId,
				InstrumentalistId: reserves[i].InstrumentalistId,
				EquipmentsId:      equipmentsId,
				TimeInterval:      reserves[i].TimeInterval,
			})
		}

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, UserReservesResponse{
			Reserves: reservesExt,
		})
	}
}

func GetStudioHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetStudioHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		//controller.ErrorResponse(wrappedWriter, fmt.Errorf("%d", idInt).Error(), http.StatusBadRequest)
		studio, err := app.StudioSvc.Get(
			&dto.GetStudioRequest{
				Id: int64(idInt),
			},
		)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, GetStudioResponse{
			Studio: studio,
		})
	}
}

func GetRoomHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetRoomsHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		//controller.ErrorResponse(wrappedWriter, fmt.Errorf("%d", idInt).Error(), http.StatusBadRequest)
		room, err := app.RoomSvc.Get(
			&dto.GetRoomRequest{
				Id: int64(idInt),
			},
		)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, GetRoomResponse{
			Room: room,
		})
	}
}

func GetProducerHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetProducerHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		//controller.ErrorResponse(wrappedWriter, fmt.Errorf("%d", idInt).Error(), http.StatusBadRequest)
		producer, err := app.ProducerSvc.Get(
			&dto.GetProducerRequest{
				Id: int64(idInt),
			},
		)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, GetProducerResponse{
			Producer: producer,
		})

	}
}

func GetInstrumentalistHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetInstrumentalistHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		//controller.ErrorResponse(wrappedWriter, fmt.Errorf("%d", idInt).Error(), http.StatusBadRequest)
		instrumentalist, err := app.InstrumentalistSvc.Get(
			&dto.GetInstrumentalistRequest{
				Id: int64(idInt),
			},
		)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, GetInstrumentalistResponse{
			Instrumentalist: instrumentalist,
		})

	}
}

func GetEquipmentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetEquipmentHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		//controller.ErrorResponse(wrappedWriter, fmt.Errorf("%d", idInt).Error(), http.StatusBadRequest)
		equipment, err := app.EquipmentSvc.Get(
			&dto.GetEquipmentRequest{
				Id: int64(idInt),
			},
		)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, GetEquipmentResponse{
			Equipment: equipment,
		})

	}
}

func AddStudioHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "AddStudioHandler"

		var req AddStudioRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.StudioSvc.Add(r.Context(), &dto.AddStudioRequest{
			Name: req.Name,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func AddRoomHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "AddRoomHandler"

		var req AddRoomRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.RoomSvc.Add(r.Context(), &dto.AddRoomRequest{
			Name:      req.Name,
			StudioId:  req.StudioId,
			StartHour: req.StartHour,
			EndHour:   req.EndHour,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func AddProducerHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "AddProducerHandler"

		var req AddProducerRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.ProducerSvc.Add(r.Context(), &dto.AddProducerRequest{
			Name:      req.Name,
			StudioId:  req.StudioId,
			StartHour: req.StartHour,
			EndHour:   req.EndHour,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func AddInstrumentalistHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "AddInstrumentalistHandler"

		var req AddInstrumentalistRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.InstrumentalistSvc.Add(r.Context(), &dto.AddInstrumentalistRequest{
			Name:      req.Name,
			StudioId:  req.StudioId,
			StartHour: req.StartHour,
			EndHour:   req.EndHour,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func AddEquipmentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "AddEquipmentHandler"

		var req AddEquipmentRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.EquipmentSvc.Add(r.Context(), &dto.AddEquipmentRequest{
			Name:     req.Name,
			StudioId: req.StudioId,
			Type:     req.TypeId,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateStudioHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateStudioHandler"

		var req UpdateStudioRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.StudioSvc.Update(&dto.UpdateStudioRequest{
			Id:   int64(idInt),
			Name: req.Name,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateRoomHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateRoomHandler"

		var req UpdateRoomRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.RoomSvc.Update(&dto.UpdateRoomRequest{
			Id:        int64(idInt),
			Name:      req.Name,
			StudioId:  req.StudioId,
			StartHour: req.StartHour,
			EndHour:   req.EndHour,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateProducerHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateProducerHandler"

		var req UpdateProducerRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.ProducerSvc.Update(&dto.UpdateProducerRequest{
			Id:        int64(idInt),
			Name:      req.Name,
			StudioId:  req.StudioId,
			StartHour: req.StartHour,
			EndHour:   req.EndHour,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateInstrumentalistHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateInstrumentalistHandler"

		var req UpdateInstrumentalistRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.InstrumentalistSvc.Update(&dto.UpdateInstrumentalistRequest{
			Id:        int64(idInt),
			Name:      req.Name,
			StudioId:  req.StudioId,
			StartHour: req.StartHour,
			EndHour:   req.EndHour,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func UpdateEquipmentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UpdateEquipmentHandler"

		var req UpdateEquipmentRequest

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.EquipmentSvc.Update(&dto.UpdateEquipmentRequest{
			Id:       int64(idInt),
			Name:     req.Name,
			StudioId: req.StudioId,
			Type:     req.TypeId,
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusUnauthorized)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func GetRoomsByStudioHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetRoomsByStudioHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id студии", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id студии к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		//controller.ErrorResponse(wrappedWriter, fmt.Errorf("%d", idInt).Error(), http.StatusBadRequest)
		rooms, err := app.RoomSvc.GetByStudio(
			&dto.GetRoomByStudioRequest{
				StudioId: int64(idInt),
			},
		)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, GetRoomsByStudioResponse{
			Rooms: rooms,
		})

	}
}

func GetInstrumentalistsByStudioHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetInstrumentalistsByStudioHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id студии", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id студии к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		//controller.ErrorResponse(wrappedWriter, fmt.Errorf("%d", idInt).Error(), http.StatusBadRequest)
		instrumentalists, err := app.InstrumentalistSvc.GetByStudio(
			&dto.GetInstrumentalistByStudioRequest{
				StudioId: int64(idInt),
			},
		)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, GetInstrumentalistsByStudioResponse{
			Instrumentalists: instrumentalists,
		})

	}
}

func GetEquipmentsByStudioHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetEquipmentsByStudioHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id студии", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id студии к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		//controller.ErrorResponse(wrappedWriter, fmt.Errorf("%d", idInt).Error(), http.StatusBadRequest)
		equipments, err := app.EquipmentSvc.GetByStudio(
			&dto.GetEquipmentByStudioRequest{
				StudioId: int64(idInt),
			},
		)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, GetEquipmentsByStudioResponse{
			Equipments: equipments,
		})

	}
}

func DeleteStudioHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UserStudioHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id студии", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id студии к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.StudioSvc.Delete(
			&dto.DeleteStudioRequest{
				Id: int64(idInt),
			},
		)

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusNoContent, nil)
	}
}

func DeleteRoomHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UserReservesHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id комнаты", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id комнаты к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.RoomSvc.Delete(
			&dto.DeleteRoomRequest{
				Id: int64(idInt),
			},
		)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusNoContent, nil)
	}
}

func DeleteProducerHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UserReservesHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id продюсера", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id продюсера к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.ProducerSvc.Delete(
			&dto.DeleteProducerRequest{
				Id: int64(idInt),
			},
		)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusNoContent, nil)
	}
}

func DeleteInstrumentalistHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UserReservesHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id инструменталиста", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id инструменталиста к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.InstrumentalistSvc.Delete(
			&dto.DeleteInstrumentalistRequest{
				Id: int64(idInt),
			},
		)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusNoContent, nil)
	}
}

func DeleteEquipmentHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "UserEquipmentHandler"
		start := time.Now()

        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =

		id := chi.URLParam(r, "id")
		if id == "" {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: пустой id оборудования", prompt).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id оборудования к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.EquipmentSvc.Delete(
			&dto.DeleteEquipmentRequest{
				Id: int64(idInt),
			},
		)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusNoContent, nil)
	}
}

func AddPhotoHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "AddPhotoHandler"

		const multiFormSizeDefault = 10000000

		start := time.Now()

        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()

		id, err := controller.GetStringClaimFromJWT(r.Context(), "id")
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = r.ParseMultipartForm(multiFormSizeDefault) // Limit your max memory usage
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		f, _, err := r.FormFile("image")
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()

		// photo, err := f.Open()
		// if err != nil {
		// 	controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
		// 	return
		// }
		// defer photo.Close()
	
		photoData, err := io.ReadAll(f)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}
	
		err = app.PhotoSvc.CreatePhoto(r.Context(), &dto.CreatePhotoRequest{
			UserId: int64(idInt),
			Data:       photoData,
		})

		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusInternalServerError)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)
	}
}

func GetPhotoHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "GetPhotoHandler"

		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()


		id, err := controller.GetStringClaimFromJWT(r.Context(), "id")
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}
		//err := json.NewDecoder(r.Body).Decode(&req)
		//if err != nil {
		//	controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
		//	return
		//}

		//ua =
		//controller.ErrorResponse(wrappedWriter, fmt.Errorf("%d", idInt).Error(), http.StatusBadRequest)

		photoData, err := app.PhotoSvc.GetPhoto(r.Context(), &dto.GetPhotoRequest{
			UserId: int64(idInt),
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, photoData)
	}
}

func DeletePhotoHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prompt := "DeletePhotoHandler"
	
		start := time.Now()
        wrappedWriter := &controller.StatusResponseWriter{ResponseWriter: w, StatusCodeOuter: http.StatusOK}

        defer func() {
    	        controller.ObserveRequest(time.Since(start), wrappedWriter.StatusCode(), r.Method, prompt)
        }()


		id, err := controller.GetStringClaimFromJWT(r.Context(), "id")
		if err != nil {
			controller.ErrorResponse(w, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: преобразование id к int: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		err = app.PhotoSvc.DeletePhoto(r.Context(), &dto.DeletePhotoRequest{
			UserId: int64(idInt),
		})
		if err != nil {
			controller.ErrorResponse(wrappedWriter, fmt.Errorf("%s: %w", prompt, err).Error(), http.StatusBadRequest)
			return
		}

		controller.SuccessResponse(wrappedWriter, http.StatusOK, nil)

	}
}
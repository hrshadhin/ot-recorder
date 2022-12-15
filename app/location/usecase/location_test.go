package usecase_test

import (
	"context"
	"errors"
	"ot-recorder/app/location/usecase"
	"ot-recorder/app/model"
	"ot-recorder/app/model/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLocation(t *testing.T) {
	mockLocationRepo := new(mocks.LocationRepository)

	mockLocation := model.Location{
		Username:  "dev",
		Device:    "phoneAndroid",
		CreatedAt: time.Now().Unix(),
		Acc:       13,
		Alt:       -42,
		Batt:      40,
		Bs:        1,
		Lat:       23.0000000,
		Lon:       90.0000000,
		M:         1,
		Tid:       "p1",
		Vac:       1,
		Vel:       0,
		Bssid:     "c0:00:00:00:00:00",
		Ssid:      "dev-test",
		IP:        "127.0.0.1",
	}

	t.Run("success", func(t *testing.T) {
		tMockLoc := mockLocation
		tMockLoc.ID = 1

		mockLocationRepo.On("CreateLocation", mock.Anything, mock.AnythingOfType("*model.Location")).
			Return(nil).Once()

		u := usecase.NewLocationUsecase(mockLocationRepo, time.Second*2)

		err := u.Ping(context.TODO(), &tMockLoc)
		assert.NoError(t, err)
		mockLocationRepo.AssertExpectations(t)
	})
}

func TestGetUserLastLocation(t *testing.T) {
	mockLocationRepo := new(mocks.LocationRepository)
	mockLocation := model.Location{
		Username:  "dev",
		Device:    "phoneAndroid",
		CreatedAt: time.Now().Unix(),
		Acc:       13,
		Alt:       -42,
		Batt:      40,
		Bs:        1,
		Lat:       23.0000000,
		Lon:       90.0000000,
		M:         1,
		Tid:       "p1",
		Vac:       1,
		Vel:       0,
		Bssid:     "c0:00:00:00:00:00",
		Ssid:      "dev-test",
		IP:        "127.0.0.1",
	}

	t.Run("success", func(t *testing.T) {
		existingLocation := mockLocation
		existingLocation.ID = 2

		mockLocationRepo.On("GetUserLastLocation", mock.Anything, mock.AnythingOfType("string")).
			Return(existingLocation, nil).Once()

		u := usecase.NewLocationUsecase(mockLocationRepo, time.Second*2)
		details, err := u.LastLocation(context.TODO(), "dev")

		assert.NoError(t, err)
		assert.Equal(t, existingLocation.Username, details.Username)
		assert.Equal(t, existingLocation.IP, details.IPAddress)
		mockLocationRepo.AssertExpectations(t)
	})

	t.Run("not-found", func(t *testing.T) {
		mockLocationRepo.On("GetUserLastLocation", mock.Anything, mock.AnythingOfType("string")).
			Return(model.Location{}, errors.New("no row found")).Once()

		u := usecase.NewLocationUsecase(mockLocationRepo, time.Second*2)
		_, err := u.LastLocation(context.TODO(), "none")

		assert.Error(t, err)
		mockLocationRepo.AssertExpectations(t)
	})
}

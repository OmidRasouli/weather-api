package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/OmidRasouli/weather-api/internal/application/service"
	"github.com/OmidRasouli/weather-api/internal/application/service/mocks"
	"github.com/OmidRasouli/weather-api/internal/domain/weather"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWeatherService_GetWeatherByID(t *testing.T) {
	type fields struct {
		repo  *mocks.MockWeatherRepository
		api   *mocks.MockAPIClient
		cache *mocks.MockCache
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	now := time.Now()
	tests := []struct {
		name        string
		setupMock   func(f fields, a args)
		args        args
		wantErr     bool
		expected    *weather.Weather
		expectedErr string
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				id:  uuid.New(),
			},
			setupMock: func(f fields, a args) {
				expected := &weather.Weather{
					ID:          a.id,
					City:        "tehran",
					Country:     "IR",
					Temperature: 30.5,
					Description: "sunny",
					Humidity:    40,
					WindSpeed:   5.5,
					FetchedAt:   now,
				}
				// Cache miss (return an error to indicate miss), then repo hit, then cache set
				f.cache.On("Get", a.ctx, a.id.String(), mock.Anything).Return(fmt.Errorf("cache miss"))
				f.repo.On("FindByID", a.ctx, mock.MatchedBy(func(u string) bool {
					return u == a.id.String()
				})).Return(expected, nil)
				f.cache.On("Set", a.ctx, a.id.String(), mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
			expected: &weather.Weather{
				ID:          uuid.Nil, // will check fields below
				City:        "tehran",
				Country:     "IR",
				Temperature: 30.5,
				Description: "sunny",
				Humidity:    40,
				WindSpeed:   5.5,
				FetchedAt:   now,
			},
		},
		{
			name: "not found",
			args: args{
				ctx: context.TODO(),
				id:  uuid.New(),
			},
			setupMock: func(f fields, a args) {
				// Cache miss (return an error), then repo not found
				f.cache.On("Get", a.ctx, a.id.String(), mock.Anything).Return(fmt.Errorf("cache miss"))
				f.repo.On("FindByID", a.ctx, mock.MatchedBy(func(u string) bool {
					return u == a.id.String()
				})).Return((*weather.Weather)(nil), fmt.Errorf("not found"))
			},
			wantErr:     true,
			expected:    nil,
			expectedErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockWeatherRepository)
			api := new(mocks.MockAPIClient)
			cache := new(mocks.MockCache)
			f := fields{repo, api, cache}
			if tt.setupMock != nil {
				tt.setupMock(f, tt.args)
			}
			service := service.NewWeatherService(repo, api, cache)
			got, err := service.GetWeatherByID(tt.args.ctx, tt.args.id.String())
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.id, got.ID)
				assert.Equal(t, tt.expected.City, got.City)
				assert.Equal(t, tt.expected.Country, got.Country)
				assert.Equal(t, tt.expected.Temperature, got.Temperature)
				assert.Equal(t, tt.expected.Description, got.Description)
				assert.Equal(t, tt.expected.Humidity, got.Humidity)
				assert.Equal(t, tt.expected.WindSpeed, got.WindSpeed)
			}
			repo.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}

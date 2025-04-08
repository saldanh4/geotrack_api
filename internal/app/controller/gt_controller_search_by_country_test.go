package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	e "geotrack_api/config/customerrors"
	l "geotrack_api/config/logger"
	mocks "geotrack_api/config/test/mocks"
	m "geotrack_api/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

var countryResult []m.GeoLocationData = []m.GeoLocationData{
	{
		Id:               1,
		As:               "Google LLC",
		City:             "Mountain View",
		Country:          "United States",
		CountryCode:      "US",
		Isp:              "Google LLC",
		Lat:              37.38605,
		Lon:              -122.08385,
		Org:              "Google LLC",
		Query:            "8.8.8.8",
		Region:           "California",
		RegionName:       "California",
		Status:           "success",
		Timezone:         "America/Los_Angeles",
		Zip:              "94035",
		TimeStamp:        "2025-04-05T20:20:00",
		DistanceSeSquare: 0.0,
	},
}

type GetCountryTestCase struct {
	name          string
	urlTarget     string
	inputControl  string
	givenCountry  *m.GivenData
	expctedResult []m.GeoLocationData
	body          string
	expctedStatus int
	entryError    *e.CustomError
	usecaseError  *e.CustomError
}

func TestGetEntriesByCountry(t *testing.T) {
	l.LoggerInit()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []GetCountryTestCase{
		{
			name:          "Valid Country",
			urlTarget:     "/?country=usa",
			inputControl:  "country",
			givenCountry:  &m.GivenData{Country: "usa"},
			expctedResult: countryResult,
			body:          "",
			expctedStatus: http.StatusOK,
			entryError:    nil,
			usecaseError:  nil,
		},
		{
			name:          "Invalid Country",
			urlTarget:     "/?country=123",
			inputControl:  "country",
			givenCountry:  &m.GivenData{Country: "123"},
			expctedResult: nil,
			body:          "",
			expctedStatus: http.StatusBadRequest,
			entryError:    e.CustomErr(e.ErrInvalidInput, "nome ou código de país invalido"),
			usecaseError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := mocks.NewMockCheckService(ctrl)
			mockUsecase := mocks.NewMockGeotrackUsecase(ctrl)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodGet, tc.urlTarget, bytes.NewBufferString(tc.body))

			mockService.EXPECT().CheckEntryData(tc.inputControl, gomock.Any()).Return(tc.givenCountry, tc.entryError)

			if tc.entryError == nil {
				mockUsecase.EXPECT().GetEntriesByCountryUsecase(tc.givenCountry.Country).Return(&tc.expctedResult, tc.usecaseError)
			}

			ipController := &GeotrackControllerImpl{
				GeotrackUsecase: mockUsecase,
				Service:         mockService,
			}
			ipController.GetEntriesByCountry(c)

			assert.Equal(t, tc.expctedStatus, recorder.Code)
			if tc.expctedStatus == http.StatusOK {
				var responseBody []m.GeoLocationData
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expctedResult, responseBody)
			} else {
				assert.Contains(t, recorder.Body.String(), tc.entryError.CustomMsg)
			}
		})
	}
}

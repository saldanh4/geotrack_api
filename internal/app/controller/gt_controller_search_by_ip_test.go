package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	e "geotrack_api/config/customerrors"
	mocks "geotrack_api/config/test/mocks"
	m "geotrack_api/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

type IpDataTest struct {
	Ip      string
	Isp     string
	Country string
	Count   int8
}

type GetTestCase struct {
	name           string
	urlTarget      string
	inputControl   string
	givenIp        *m.GivenData
	expectedResult *m.GeoLocationData
	body           string
	expectedStatus int
	entryError     *e.CustomError
	usecaseError   *e.CustomError
	ipData         IpDataTest
}

func TestGetEntriesByIp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCase := []GetTestCase{
		{
			name:           "Valid IP",
			urlTarget:      "/?ip=8.8.8.8",
			inputControl:   "ip",
			givenIp:        &m.GivenData{Ip: "8.8.8.8"},
			expectedResult: &m.GeoLocationData{Query: "8.8.8.8", Isp: "Google LLC", Country: "United States", Count: 1},
			body:           "",
			expectedStatus: 200,
			entryError:     nil,
			usecaseError:   nil,
			ipData:         IpDataTest{Ip: "8.8.8.8", Isp: "Google LLC", Country: "United States", Count: 1},
		},
		{
			name:           "Invalid IP",
			urlTarget:      "/?ip=Teste",
			inputControl:   "ip",
			givenIp:        &m.GivenData{Ip: "Teste"},
			expectedResult: nil,
			body:           "",
			expectedStatus: 400,
			entryError:     e.CustomErr(e.ErrInvalidInput, "formato de ip inválido"),
			usecaseError:   nil,
			ipData:         IpDataTest{},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := mocks.NewMockGeotrackUsecase(ctrl)
			mockService := mocks.NewMockCheckService(ctrl)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodGet, tc.urlTarget, bytes.NewBufferString(tc.body))

			mockService.EXPECT().CheckEntryData(tc.inputControl, gomock.Any()).Return(tc.givenIp, tc.entryError)

			if tc.entryError == nil {
				mockUsecase.EXPECT().GetEntriesByIpUsecase(tc.givenIp.Ip).Return(tc.expectedResult, tc.usecaseError)
			}

			ipController := &GeotrackControllerImpl{
				GeotrackUsecase: mockUsecase,
				Service:         mockService,
			}
			ipController.GetEntriesByIp(c)

			var response IpDataTest
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.ipData, response)

			assert.Equal(t, tc.expectedStatus, recorder.Code)

		})

	}
}

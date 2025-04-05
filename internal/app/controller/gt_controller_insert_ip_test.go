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

type TesteCase struct {
	name             string
	urlTarget        string
	body             string
	inputControl     string
	givenIp          *m.GivenData
	expectedResult   *m.GeoLocationData
	entryError       *e.CustomError
	useCaseError     *e.CustomError
	expectedStatus   int
	expectedResponse gin.H
}

//Ajustar as interfaces para separa os serviços CheckEntryData
//Ajustar as assinaturas das chamadas
//Refazer os mocks (mockgen)

func TestCreateIP(t *testing.T) {
	l.LoggerInit()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGivenIP := &m.GivenData{
		Ip: "8.8.8.8",
	}

	mockExpectedResult := &m.GeoLocationData{
		Query:     "8.8.8.8",
		Country:   "USA",
		TimeStamp: "2025-03-14",
	}

	mockEntryError := &e.CustomError{
		BaseError: e.ErrInvalidInput,
		CustomMsg: "solicitações POST não devem ter dados enviados via url",
	}

	mockUsecaseError := &e.CustomError{
		BaseError: e.ErrInternalServer,
		CustomMsg: "não foram encontrados dados de geolocalização para o IP informado",
	}

	testCases := []TesteCase{
		{
			"Success - Valid Input",
			"/add_ip",
			`{"ip":"8.8.8.8"}`,
			"ip",
			mockGivenIP,
			mockExpectedResult,
			nil,
			nil,
			200,
			gin.H{
				"ip":         "8.8.8.8",
				"country":    "USA",
				"insertDate": "2025-03-14",
			},
		},
		{
			"Fail - Invalid Input - POST With URL",
			"/add_ip?ip=8.8.8.8",
			"",
			"ip",
			mockGivenIP,
			mockExpectedResult,
			mockEntryError,
			nil,
			400,
			gin.H{
				"message": "solicitações POST não devem ter dados enviados via url",
			},
		},
		{
			"Fail - Internal Error",
			"/add_ip",
			`{"ip":"8.8.8.8"}`,
			"ip",
			mockGivenIP,
			mockExpectedResult,
			nil,
			mockUsecaseError,
			500,
			gin.H{
				"message": "não foram encontrados dados de geolocalização para o IP informado",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			//mockController := mocks.NewMockGeotrackController(ctrl)
			mockUsecase := mocks.NewMockGeotrackUsecase(ctrl)
			mockService := mocks.NewMockCheckService(ctrl)

			//mockController := mocks.NewMockGeotrackController(ctrl)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodPost, tc.urlTarget, bytes.NewBufferString(tc.body))
			c.Request.Header.Set("Content-Type", "application/json")

			mockService.EXPECT().CheckEntryData(tc.inputControl, gomock.Any()).Return(tc.givenIp, tc.entryError)

			if tc.entryError == nil {
				mockUsecase.EXPECT().CreateIP(tc.givenIp).Return(tc.expectedResult, tc.useCaseError)
			}

			ipController := &GeotrackControllerImpl{
				GeotrackUsecase: mockUsecase,
				Service:         mockService,
			}
			ipController.CreateIP(c)

			assert.Equal(t, tc.expectedStatus, recorder.Code)

			if tc.expectedResponse != nil {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			}

		})
	}
}

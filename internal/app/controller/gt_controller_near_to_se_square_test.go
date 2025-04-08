package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type TestResult struct {
	Ip       string
	As       string
	City     string
	Country  string
	Distance string
}

type SeSqrTestCase struct {
	name       string
	urlTarget  string
	queryParam string
	body       string
	geoData    m.GeoLocationData
	result     TestResult
	usecaseErr *e.CustomError
	expected   gin.H
}

func TestGetNearestIpToSeSquare(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCase := []SeSqrTestCase{
		{
			name:       "Valid return",
			urlTarget:  "/nearest_se_square",
			queryParam: "",
			body:       "",
			geoData: m.GeoLocationData{
				Query:   "8.8.8.8",
				As:      "AS15169 Google LLC",
				City:    "Ashburn",
				Country: "United States",
			},
			result:     TestResult{},
			usecaseErr: nil,
			expected:   gin.H{},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := mocks.NewMockGeotrackUsecase(ctrl)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodGet, tc.urlTarget+tc.queryParam, bytes.NewBufferString(tc.body))
			c.Request.Header.Set("Content-Type", "application/json")

			mockUsecase.EXPECT().GetNearestIpToSeSquareUsecase().Return(&tc.geoData, tc.usecaseErr)

			ipController := &GeotrackControllerImpl{
				GeotrackUsecase: mockUsecase,
				Service:         NewDefaultCheckService(),
			}
			ipController.GetNearestIpToSeSquare(c)
			tc.result = TestResult{
				Ip:       tc.geoData.Query,
				As:       tc.geoData.As,
				City:     tc.geoData.City,
				Country:  tc.geoData.Country,
				Distance: fmt.Sprintf("%.2fkm", tc.geoData.DistanceSeSquare),
			}
			tc.expected = gin.H{
				"Ip mais próximo da praça da Sé é": tc.result.Ip,
				"dados":                            tc.result}

			assert.Equal(t, http.StatusOK, recorder.Code)
			var responseBody map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &responseBody)

			dados := responseBody["dados"].(map[string]interface{})
			assert.Equal(t, tc.result.Ip, dados["Ip"])
			assert.Equal(t, tc.result.As, dados["As"])
			assert.Equal(t, tc.result.City, dados["City"])
			assert.Equal(t, tc.result.Country, dados["Country"])
			assert.Equal(t, tc.result.Distance, dados["Distance"])
		})
	}
}

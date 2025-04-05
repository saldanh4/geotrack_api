package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	e "geotrack_api/config/customerrors"
	mocks "geotrack_api/config/test/mocks"
	m "geotrack_api/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type TestCase struct {
	name         string
	method       string
	inputControl string
	urlQuery     string
	body         string
	result       *m.GivenData
	checkErr     *e.CustomError
	validateErr  *e.CustomError
}

func TestCheckEntryData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockGeotrackUsecase(ctrl)
	mockCheckData := mocks.NewMockCheckData(ctrl)

	testCase := []TestCase{
		{
			name:         "GET method - Valid IP",
			method:       http.MethodGet,
			inputControl: "ip",
			urlQuery:     "/?ip=8.8.8.8",
			body:         "",
			result:       &m.GivenData{Ip: "8.8.8.8"},
			checkErr:     nil,
			validateErr:  nil,
		},
		// {
		// 	"GET method - ",
		// 	http.MethodGet,
		// 	"teste",
		// 	"/?ip=8.8.8.8",
		// 	"",
		// 	nil, //&m.GivenData{Ip: "8.8.8.8"},
		// 	e.CustomErr(e.ErrInternalServer, "falha ao checar se a entrada é 'ip' ou 'país'"),
		// 	nil,
		// },
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			//mockService := mocks.NewMockCheckService(ctrl)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(tc.method, tc.urlQuery, bytes.NewBufferString(tc.body))

			//t.Logf("Requisição criada: Método = %s, URL = %s, Body = %s", c.Request.Method, c.Request.URL.String(), c.Request.Body)

			//mockCheckData.EXPECT().CheckInputData(tc.inputControl, gomock.Any()).Return(tc.result, nil)
			mockCheckData.EXPECT().CheckInputData("ip", gomock.Any()).Return(&m.GivenData{Ip: "8.8.8.8"}, nil)
			if tc.checkErr == nil {
				switch tc.inputControl {
				case "ip":
					//mockCheckData.EXPECT().ValidateIp(tc.result.Ip).Return(nil)
					mockCheckData.EXPECT().ValidateIp("8.8.8.8").Return(nil)
					// t.Logf("Validando IP: %s", tc.result.Ip)
					// t.Logf("Verificando erro: %s", tc.validateErr)
				case "country":
					mockCheckData.EXPECT().ValidateCountry(tc.result.Country).Return(nil)
					// t.Logf("Validando IP: %s", tc.result.Ip)
					// t.Logf("Verificando erro: %s", tc.validateErr)
				}
			}

			controller := &GeotrackControllerImpl{
				GeotrackUsecase: mockUsecase,
				Service: &CheckServiceImpl{
					CheckData: mockCheckData,
				},
			}
			result, err := controller.Service.CheckEntryData(tc.inputControl, c)
			t.Logf("Erro retornado pela função CheckEntryData: %v", err)
			t.Logf("Valor para result = %v", result)
			v := err == tc.validateErr
			t.Logf("Erro é igual Validate Erro? %v", v)

			if err != nil { //tc.checkErr != nil || tc.validateErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				t.Logf("erro = %v", err)
				//assert.NoError(t, err)

				assert.NotNil(t, result)
				assert.Equal(t, tc.result.Ip, result.Ip)
			}
		})

	}

}

// func TestCheckEntry(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	// 	defer ctrl.Finish()
// 	//mockUsecase := new(mocks.MockGeotrackUsecase)
// 	mockService := mocks.NewMockCheckService(ctrl)
// 	//controller := NewGeotrackController(mockUsecase, mockService)

// 	testCases := []struct {
// 		name         string
// 		method       string
// 		inputControl string
// 		query        string
// 		body         string
// 		hasError     bool
// 		expected     *m.GivenData
// 		tError       *e.CustomError
// 	}{
// 		//Casos de teste GET IP
// 		{"GET method - Valid IP", http.MethodGet, "ip", "/?ip=8.8.8.8", "", false, &m.GivenData{Ip: "8.8.8.8"}, nil},
// 		{"GET method - Invalid Input Control", http.MethodGet, "teste", "/?ip=8.8.8.8", "", true, &m.GivenData{Ip: "8.8.8.8"}, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"GET Method - Invalid IP", http.MethodGet, "ip", "/?ip=256.256.256.256", "", true, &m.GivenData{Ip: "256.256.256.256"}, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"GET Method - Invalid IP", http.MethodGet, "ip", "/?ip=::1", "", true, &m.GivenData{Ip: "::1"}, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"GET Method - Invalid IP", http.MethodGet, "ip", "/?ip=", "", true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"GET Method - With body", http.MethodGet, "ip", "/?ip=8.8.8.8", `{"ip":"8.8.8.8"}`, true, nil, e.CustomErr(e.ErrInvalidInput, "")},

// 		//Casos de teste GET Country
// 		{"GET method - Valid Country", http.MethodGet, "country", "/?country=Brazil", "", false, &m.GivenData{Country: "Brazil"}, nil},
// 		{"GET method - Invalid Country", http.MethodGet, "country", "/?country=134", "", true, &m.GivenData{Country: "134"}, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"GET method - Invalid Country", http.MethodGet, "country", "/?country=", "", true, nil, e.CustomErr(e.ErrInvalidInput, "")},

// 		//Casos de teste POST IP
// 		{"POST method - Valid IP", http.MethodPost, "ip", "/", `{"ip": "8.8.8.8"}`, false, &m.GivenData{Ip: "8.8.8.8"}, nil},
// 		{"POST method - Invalid Input Control", http.MethodPost, "teste", "/", `{"ip": "8.8.8.8"}`, true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"POST method - Invalid IP", http.MethodPost, "ip", "/", `{"ip": "256.256.256"}`, true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"POST method - Invalid IP", http.MethodPost, "ip", "/", `{"ip": "::1"}`, true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"POST method - with URL", http.MethodPost, "ip", "/?ip=8.8.8.8", `{"ip": "8.8.8.8"}`, true, nil, e.CustomErr(e.ErrInvalidInput, "")},

// 		//Casos de test DELETE
// 		{"DELETE Method - Valid IP", http.MethodDelete, "ip", "/?ip=8.8.8.8", "", false, &m.GivenData{Ip: "8.8.8.8"}, nil},
// 		{"DELETE Method - Invalid Input Control", http.MethodDelete, "", "/?ip=8.8.8.8", "", true, &m.GivenData{Ip: "8.8.8.8"}, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"DELETE Method - Invalid IP", http.MethodDelete, "ip", "/?ip=256.256.256.256", "", true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"DELETE Method - Invalid IP", http.MethodDelete, "ip", "/?ip=::1", "", true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"DELETE Method - Invalid IP", http.MethodDelete, "ip", "/?ip=", "", true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"DELETE Method - Valid Country", http.MethodDelete, "country", "/?country=Brazil", "", false, &m.GivenData{Country: "Brazil"}, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"DELETE Method - Invalid Country", http.MethodDelete, "country", "/?country=B", "", true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"DELETE Method - Invalid Country", http.MethodDelete, "country", "/?country=", "", true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 		{"DELETE Method - With body", http.MethodDelete, "ip", "/?ip=8.8.8.8", `{"ip":"8.8.8.8"}`, true, nil, e.CustomErr(e.ErrInvalidInput, "")},

// 		//Caso de teste Método inválido
// 		{"Invalid Method", http.MethodPut, "country", "/?country=B", "", true, nil, e.CustomErr(e.ErrInvalidInput, "")},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			gin.SetMode(gin.TestMode)
// 			recorder := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(recorder)
// 			c.Request = httptest.NewRequest(tc.method, tc.query, bytes.NewBufferString(tc.body))

// 			mockService.EXPECT().CheckEntryData(tc.inputControl, gomock.Any()).Return(tc.expected, tc.tError)
// 			//result, err := mockService.CheckEntryData(tc.inputControl, c)

// 			if tc.hasError {
// 				assert.Nil(t, tc.expected)
// 				assert.NotNil(t, tc.tError)
// 			} else {
// 				assert.NotNil(t, tc.expected)
// 				assert.Nil(t, tc.tError)
// 				if tc.inputControl == "ip" {
// 					assert.Equal(t, tc.expected.Ip, tc.expected.Ip)
// 				} else if tc.inputControl == "country" {
// 					assert.Equal(t, tc.expected.Country, tc.expected.Country)
// 				}
// 			}
// 		})
// 	}
// }

// func SubFunctionsTest() {
// 	// CheckInputData(input string, c *gin.Context) (*m.GivenData, *e.CustomError)
// 	// ValidateIp(ip string) *customerrors.CustomError
// 	// ValidateCountry(country string) *customerrors.CustomError

// }

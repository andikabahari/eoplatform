package usecase

import (
	"net/http"
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	mr "github.com/andikabahari/eoplatform/repository/mock_repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type bankAccountUsecaseSuite struct {
	suite.Suite

	ctrl                  *gomock.Controller
	bankAccountRepository *mr.MockBankAccountRepository

	usecase BankAccountUsecase
}

func (s *bankAccountUsecaseSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.bankAccountRepository = mr.NewMockBankAccountRepository(s.ctrl)

	s.usecase = NewBankAccountUsecase(s.bankAccountRepository)
}

func (s *bankAccountUsecaseSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestBankAccountUsecaseSuite(t *testing.T) {
	suite.Run(t, new(bankAccountUsecaseSuite))
}

func (s *bankAccountUsecaseSuite) TestGetBankAccount() {
	testCases := []struct {
		Name         string
		Body         any
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"ok",
			nil,
			&helper.JWTCustomClaims{ID: 1, Role: "organizer"},
			func() {
				s.bankAccountRepository.EXPECT().FindByUserID(
					gomock.Eq(&model.BankAccount{}),
					gomock.Eq(uint(1)),
				)
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			s.usecase.GetBankAccount(testCase.Claims, &model.BankAccount{})
		})
	}
}

func (s *bankAccountUsecaseSuite) TestCreateBankAccount() {
	testCases := []struct {
		Name         string
		Body         *request.CreateBankAccountRequest
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"ok",
			&request.CreateBankAccountRequest{
				BasicBankAccount: request.BasicBankAccount{
					Bank:     "bni",
					VANumber: "12345",
				},
			},
			&helper.JWTCustomClaims{ID: 1, Role: "organizer"},
			func() {
				s.bankAccountRepository.EXPECT().Create(
					gomock.Eq(&model.BankAccount{
						Bank:     "bni",
						VANumber: "12345",
						UserID:   1,
					}),
				)
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			s.usecase.CreateBankAccount(testCase.Claims, &model.BankAccount{}, testCase.Body)
		})
	}
}

func (s *bankAccountUsecaseSuite) TestUpdateBankAccount() {
	testCases := []struct {
		Name         string
		Body         *request.UpdateBankAccountRequest
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			&request.UpdateBankAccountRequest{
				BasicBankAccount: request.BasicBankAccount{
					Bank:     "bni",
					VANumber: "12345",
				},
			},
			&helper.JWTCustomClaims{},
			func() {
				s.bankAccountRepository.EXPECT().FindByUserID(
					gomock.Eq(&model.BankAccount{}),
					gomock.Eq(uint(0)),
				)
			},
			http.StatusNotFound,
		},
		{
			"ok",
			&request.UpdateBankAccountRequest{
				BasicBankAccount: request.BasicBankAccount{
					Bank:     "bni",
					VANumber: "12345",
				},
			},
			&helper.JWTCustomClaims{ID: 1, Role: "organizer"},
			func() {
				s.bankAccountRepository.EXPECT().FindByUserID(
					gomock.Eq(&model.BankAccount{}),
					gomock.Eq(uint(1)),
				).SetArg(0, model.BankAccount{Model: gorm.Model{ID: 1}})

				s.bankAccountRepository.EXPECT().Update(
					gomock.Eq(&model.BankAccount{Model: gorm.Model{ID: 1}}),
					gomock.Eq(&request.UpdateBankAccountRequest{
						BasicBankAccount: request.BasicBankAccount{
							Bank:     "bni",
							VANumber: "12345",
						},
					}),
				)
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.UpdateBankAccount(testCase.Claims, &model.BankAccount{}, testCase.Body); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}

package advertisementCreate

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"api-tests-template/internal/assertions"
	adClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	"api-tests-template/tests"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	defaultAdDescription = "Подробное описание тестового объявления со всеми параметрами"
	defaultAdPrice       = 3500.0
	defaultAdQuantity    = 3
	maxTitleLength       = 50
	invalidToken         = "invalid_bearer_token"
	testPhotoFile        = "test_photo.png"

	errCodeInvalidRequest = "invalid_request"
	errCodeUnauthorized   = "unauthorized"
)

type CreateAdvertisementSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestCreateAdvertisementSuite(t *testing.T) {
	suite.Run(t, new(CreateAdvertisementSuite))
}

func (s *CreateAdvertisementSuite) SetupSuite() {
	tests.SetupSuite()
	tests.Precondition("Авторизация пользователя из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *CreateAdvertisementSuite) TearDownSuite() {
	tests.TearDownSuite()
}

func (s *CreateAdvertisementSuite) TestCreateAdvertisementPositive() {
	parentT := s.T()
	tests.Precondition("Подготовка уникальных данных для создания объявления")
	uniqueTitle := fmt.Sprintf("Тестовый товар %d", time.Now().UnixNano())
	price := defaultAdPrice
	quantity := defaultAdQuantity

	params := adClient.CreateAdvertisementFormParams{
		Title:       uniqueTitle,
		Description: defaultAdDescription,
		Price:       &price,
		Quantity:    &quantity,
		PhotoPaths:  []string{testPhotoFile},
	}

	var createdAd models.Advertisement

	s.Run("Создание объявления со всеми возможными параметрами", func() {
		createdAd = advertisements.CreateAdvertisement(s.T(), s.loginData.Token, params, http.StatusCreated)
		require.NotEmpty(s.T(), createdAd.ID, "ID созданного объявления не должен быть пустым")

		parentT.Cleanup(func() {
			tests.Precondition("Очистка: удаление созданного объявления " + createdAd.ID)
			advertisements.DeleteAdvertisement(parentT, s.loginData.Token, createdAd.ID, http.StatusNoContent)
		})
	})

	s.Run("Проверка доступности и совпадения параметров через GET /advertisement", func() {
		fetchedAd := advertisements.GetAdvertisementByID(s.T(), s.loginData.Token, createdAd.ID, http.StatusOK)
		require.Equal(s.T(), createdAd.ID, fetchedAd.ID, "ID объявления не совпадает")
		require.Equal(s.T(), uniqueTitle, fetchedAd.Title, "Название объявления не совпадает")
		require.Equal(s.T(), defaultAdDescription, fetchedAd.Description, "Описание объявления не совпадает")
		require.NotNil(s.T(), fetchedAd.Price, "Цена не должна быть nil")
		require.InDelta(s.T(), price, *fetchedAd.Price, 0.001, "Цена объявления не совпадает")
		require.NotNil(s.T(), fetchedAd.Quantity, "Количество не должно быть nil")
		require.Equal(s.T(), quantity, *fetchedAd.Quantity, "Количество не совпадает")
	})

	s.Run("Проверка появления объявления в поиске GET /advertisements?search={title}", func() {
		utils.WaitUntil(s.T(), 5*time.Second, 200*time.Millisecond, "объявление должно появиться в поиске", func() bool {
			list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, uniqueTitle)
			for _, item := range list.Items {
				if item.ID == createdAd.ID && item.Title == uniqueTitle {
					return true
				}
			}
			return false
		})
	})

	s.Run("Проверка загрузки фотографии через GET /advertisements/{id}/photos", func() {
		photos := advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, createdAd.ID, http.StatusOK)
		require.Len(s.T(), photos, 1, "Должна быть загружена ровно 1 фотография")
		require.NotEmpty(s.T(), photos[0].ID, "ID фотографии не должен быть пустым")
		require.NotEmpty(s.T(), photos[0].URL, "URL загруженной фотографии не должен быть пустым")
	})
}

func (s *CreateAdvertisementSuite) TestCreateAdvertisementNegativeValidation() {
	negativePrice := -100.0
	validPrice := defaultAdPrice

	cases := []struct {
		name            string
		params          adClient.CreateAdvertisementFormParams
		expectedMessage string
	}{
		{
			name: "Пустой обязательный параметр title",
			params: adClient.CreateAdvertisementFormParams{
				Title:       "",
				Description: defaultAdDescription,
				Price:       &validPrice,
				PhotoPaths:  []string{testPhotoFile},
			},
			expectedMessage: "title is required",
		},
		{
			name: "Пустой обязательный параметр description",
			params: adClient.CreateAdvertisementFormParams{
				Title:       "Валидное название",
				Description: "",
				Price:       &validPrice,
				PhotoPaths:  []string{testPhotoFile},
			},
			expectedMessage: "description is required",
		},
		{
			name: "Отсутствует обязательная фотография",
			params: adClient.CreateAdvertisementFormParams{
				Title:       "Объявление без фото",
				Description: defaultAdDescription,
				Price:       &validPrice,
				PhotoPaths:  nil,
			},
			expectedMessage: "at least 1 photo is required",
		},
		{
			name: "Отрицательная цена",
			params: adClient.CreateAdvertisementFormParams{
				Title:       "Объявление с отрицательной ценой",
				Description: defaultAdDescription,
				Price:       &negativePrice,
				PhotoPaths:  []string{testPhotoFile},
			},
			expectedMessage: "price must be between 0 and 1000000 inclusive",
		},
		{
			name: "Слишком длинный title (больше 50 символов)",
			params: adClient.CreateAdvertisementFormParams{
				Title:       strings.Repeat("A", maxTitleLength+1),
				Description: defaultAdDescription,
				Price:       &validPrice,
				PhotoPaths:  []string{testPhotoFile},
			},
			expectedMessage: "title must not exceed 50 characters provided 51",
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			body := advertisements.CreateAdvertisementRaw(s.T(), s.loginData.Token, tc.params, http.StatusBadRequest)
			assertions.AssertErrorResponse(s.T(), body, errCodeInvalidRequest, tc.expectedMessage, "POST /advertisement ("+tc.name+")")
		})
	}
}

func (s *CreateAdvertisementSuite) TestCreateAdvertisementUnauthorized() {
	validPrice := defaultAdPrice

	params := adClient.CreateAdvertisementFormParams{
		Title:       "Объявление без авторизации",
		Description: defaultAdDescription,
		Price:       &validPrice,
		PhotoPaths:  []string{testPhotoFile},
	}

	cases := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{
			name:            "Отсутствующий токен",
			token:           "",
			expectedMessage: "Authorization header required",
		},
		{
			name:            "Некорректный токен",
			token:           invalidToken,
			expectedMessage: "Invalid or expired token",
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			body := advertisements.CreateAdvertisementRaw(s.T(), tc.token, params, http.StatusUnauthorized)
			assertions.AssertErrorResponse(s.T(), body, errCodeUnauthorized, tc.expectedMessage, "POST /advertisement ("+tc.name+")")
		})
	}
}

package favorites

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"api-tests-template/internal/assertions"
	adClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
	"api-tests-template/tests"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
)

const (
	nonExistentUUID    = "11111111-1111-4111-8111-111111111111"
	invalidBearerToken = "invalid_bearer_token"
	secondUserPassword = "TestPassword123!"
	secondUserName     = "ТестовыйПродавец"
	favAdDescription   = "Объявление для проверки функциональности избранного"
	favAdPrice         = 1990.0
)

type FavoritesSuite struct {
	suite.Suite
	mainUser   authModels.LoginOkResponse
	sellerUser authModels.LoginOkResponse
}

func TestFavoritesSuite(t *testing.T) {
	suite.Run(t, new(FavoritesSuite))
}

func (s *FavoritesSuite) SetupSuite() {
	tests.SetupSuite()
	tests.Precondition("Авторизация основного пользователя из переменных окружения")
	s.mainUser = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	tests.Precondition("Регистрация второго пользователя-продавца для создания чужих объявлений")
	sellerEmail := fmt.Sprintf("seller_%d@example.com", time.Now().UnixNano())
	s.sellerUser = auth.Register(s.T(), sellerEmail, secondUserPassword, secondUserName)
}

func (s *FavoritesSuite) TearDownSuite() {
	tests.TearDownSuite()
}

// createAdForUser создаёт тестовое объявление от имени пользователя с переданным токеном.
func (s *FavoritesSuite) createAdForUser(t *testing.T, token string) adModels.Advertisement {
	t.Helper()
	photoPath := utils.CreateTempPNGFile(t, "fav_ad.png")
	price := favAdPrice
	title := fmt.Sprintf("Товар для избранного %d", time.Now().UnixNano())

	params := adClient.CreateAdvertisementFormParams{
		Title:       title,
		Description: favAdDescription,
		Price:       &price,
		PhotoPaths:  []string{photoPath},
	}
	return advertisements.CreateAdvertisement(t, token, params, http.StatusCreated)
}

func (s *FavoritesSuite) TestFavoritesFullLifecyclePositive() {
	parentT := s.T()
	tests.Precondition("Создание объявления от имени второго пользователя")
	sellerAd := s.createAdForUser(parentT, s.sellerUser.Token)
	adDeleted := false

	parentT.Cleanup(func() {
		if !adDeleted {
			tests.Precondition("Очистка: удаление тестового объявления " + sellerAd.ID)
			advertisements.DeleteAdvertisement(parentT, s.sellerUser.Token, sellerAd.ID, http.StatusNoContent)
		}
	})

	s.Run("Добавление объявления в избранное и проверка появления в GET /favorites", func() {
		favorites.AddToFavorites(s.T(), s.mainUser.Token, sellerAd.ID, http.StatusCreated)

		favList := favorites.GetFavorites(s.T(), s.mainUser.Token, http.StatusOK)
		assertions.ContainsFavoriteByAdID(s.T(), favList.Items, sellerAd.ID)
	})

	s.Run("Удаление объявления из избранного и проверка отсутствия в списке", func() {
		favorites.RemoveFromFavorites(s.T(), s.mainUser.Token, sellerAd.ID, http.StatusNoContent)

		favList := favorites.GetFavorites(s.T(), s.mainUser.Token, http.StatusOK)
		assertions.NotContainsFavoriteByAdID(s.T(), favList.Items, sellerAd.ID)
	})

	s.Run("При удалении самого объявления через DELETE /advertisement оно исчезает из избранного", func() {
		favorites.AddToFavorites(s.T(), s.mainUser.Token, sellerAd.ID, http.StatusCreated)

		favBeforeDelete := favorites.GetFavorites(s.T(), s.mainUser.Token, http.StatusOK)
		assertions.ContainsFavoriteByAdID(s.T(), favBeforeDelete.Items, sellerAd.ID)

		advertisements.DeleteAdvertisement(s.T(), s.sellerUser.Token, sellerAd.ID, http.StatusNoContent)
		adDeleted = true

		favAfterDelete := favorites.GetFavorites(s.T(), s.mainUser.Token, http.StatusOK)
		assertions.NotContainsFavoriteByAdID(s.T(), favAfterDelete.Items, sellerAd.ID)
	})
}

func (s *FavoritesSuite) TestFavoritesNonExistentAdNegative() {
	s.Run("Добавление несуществующего advertisement_id в избранное возвращает 404", func() {
		body := favorites.AddToFavorites(s.T(), s.mainUser.Token, nonExistentUUID, http.StatusNotFound)
		require.NotEmpty(s.T(), gjson.GetBytes(body, "error").String(), "Поле error не должно быть пустым")
	})

	s.Run("Удаление несуществующего advertisement_id из избранного возвращает 404", func() {
		body := favorites.RemoveFromFavorites(s.T(), s.mainUser.Token, nonExistentUUID, http.StatusNotFound)
		require.NotEmpty(s.T(), gjson.GetBytes(body, "error").String(), "Поле error не должно быть пустым")
	})
}

func (s *FavoritesSuite) TestFavoritesBusinessConstraintsNegative() {
	parentT := s.T()
	ownAd := s.createAdForUser(parentT, s.mainUser.Token)
	sellerAd := s.createAdForUser(parentT, s.sellerUser.Token)

	parentT.Cleanup(func() {
		advertisements.DeleteAdvertisement(parentT, s.mainUser.Token, ownAd.ID, http.StatusNoContent)
		advertisements.DeleteAdvertisement(parentT, s.sellerUser.Token, sellerAd.ID, http.StatusNoContent)
	})

	s.Run("Попытка добавить собственное объявление в избранное возвращает 403 Forbidden", func() {
		body := favorites.AddToFavorites(s.T(), s.mainUser.Token, ownAd.ID, http.StatusForbidden)
		require.NotEmpty(s.T(), gjson.GetBytes(body, "error").String(), "Ожидалось сообщение об ошибке при добавлении своего объявления")
	})

	s.Run("Повторное добавление того же объявления в избранное возвращает 409 Conflict", func() {
		favorites.AddToFavorites(s.T(), s.mainUser.Token, sellerAd.ID, http.StatusCreated)
		body := favorites.AddToFavorites(s.T(), s.mainUser.Token, sellerAd.ID, http.StatusConflict)
		require.NotEmpty(s.T(), gjson.GetBytes(body, "error").String(), "Ожидалось сообщение о конфликте при повторном добавлении")
	})
}

func (s *FavoritesSuite) TestFavoritesUnauthorizedNegative() {
	cases := []struct {
		name  string
		token string
	}{
		{name: "Без токена авторизации", token: ""},
		{name: "С некорректным токеном", token: invalidBearerToken},
	}

	for _, tc := range cases {
		s.Run("GET /favorites "+tc.name, func() {
			body := favorites.GetFavoritesRaw(s.T(), tc.token, http.StatusUnauthorized)
			require.NotEmpty(s.T(), gjson.GetBytes(body, "error").String())
		})

		s.Run("POST /favorites "+tc.name, func() {
			body := favorites.AddToFavorites(s.T(), tc.token, nonExistentUUID, http.StatusUnauthorized)
			require.NotEmpty(s.T(), gjson.GetBytes(body, "error").String())
		})

		s.Run("DELETE /favorites/{ad_id} "+tc.name, func() {
			body := favorites.RemoveFromFavorites(s.T(), tc.token, nonExistentUUID, http.StatusUnauthorized)
			require.NotEmpty(s.T(), gjson.GetBytes(body, "error").String())
		})
	}
}

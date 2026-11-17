package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetFavorites получает список избранных объявлений авторизованного пользователя.
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.FavoritesPath).
		Expect(t).
		End().Response
}

// HttpPostFavorites добавляет объявление в избранное.
func HttpPostFavorites(t *testing.T, token string, body any) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(body).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного по его ID.
func HttpDeleteFavorite(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoriteByAdIDPath(adID)).
		Expect(t).
		End().Response
}

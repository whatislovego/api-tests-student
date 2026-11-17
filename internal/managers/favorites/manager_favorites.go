package favorites

import (
	"encoding/json"
	"io"
	"testing"

	favoritesClient "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"

	"github.com/stretchr/testify/require"
)

// AddToFavorites отправляет запрос на добавление объявления в избранное и проверяет статус-код.
func AddToFavorites(t *testing.T, token, adID string, expectedStatus int) []byte {
	t.Helper()
	raw := favoritesClient.HttpPostFavorites(t, token, models.AddFavoriteRequest{
		AdvertisementID: adID,
	})
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа POST /favorites")
	require.Equalf(t, expectedStatus, raw.StatusCode, "Некорректный HTTP-статус POST /favorites. Ответ: %s", string(body))

	return body
}

// RemoveFromFavorites удаляет объявление из избранного и проверяет статус-код.
func RemoveFromFavorites(t *testing.T, token, adID string, expectedStatus int) []byte {
	t.Helper()
	raw := favoritesClient.HttpDeleteFavorite(t, token, adID)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа DELETE /favorites/{ad_id}")
	require.Equalf(t, expectedStatus, raw.StatusCode, "Некорректный HTTP-статус DELETE /favorites/{ad_id}. Ответ: %s", string(body))

	return body
}

// GetFavoritesRaw запрашивает список избранного, проверяет статус-код и возвращает сырое тело ответа.
func GetFavoritesRaw(t *testing.T, token string, expectedStatus int) []byte {
	t.Helper()
	raw := favoritesClient.HttpGetFavorites(t, token)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа GET /favorites")
	require.Equalf(t, expectedStatus, raw.StatusCode, "Некорректный HTTP-статус GET /favorites. Ответ: %s", string(body))

	return body
}

// GetFavorites возвращает десериализованный список избранных объявлений пользователя.
func GetFavorites(t *testing.T, token string, expectedStatus int) models.FavoriteListResponse {
	t.Helper()
	body := GetFavoritesRaw(t, token, expectedStatus)

	var resp models.FavoriteListResponse
	require.NoError(t, json.Unmarshal(body, &resp), "Ошибка десериализации списка избранного")
	return resp
}

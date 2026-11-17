package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	advertisementsClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"

	"github.com/stretchr/testify/require"
)

// GetAdvertisementsList выполняет поиск объявлений и возвращает типизированный список.
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisementsClient.HttpGetAdvertisements(t, token, text)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// CreateAdvertisementRaw отправляет запрос на создание объявления, проверяет статус-код и возвращает тело ответа.
func CreateAdvertisementRaw(
	t *testing.T,
	token string,
	params advertisementsClient.CreateAdvertisementFormParams,
	expectedStatus int,
) []byte {
	t.Helper()
	raw := advertisementsClient.HttpPostAdvertisement(t, token, params)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа при создании объявления")
	require.Equalf(t, expectedStatus, raw.StatusCode, "Некорректный HTTP-статус при создании объявления. Ответ: %s", string(body))

	return body
}

// CreateAdvertisement создаёт объявление и возвращает распарсенную структуру Advertisement.
func CreateAdvertisement(
	t *testing.T,
	token string,
	params advertisementsClient.CreateAdvertisementFormParams,
	expectedStatus int,
) models.Advertisement {
	t.Helper()
	body := CreateAdvertisementRaw(t, token, params, expectedStatus)

	var resp models.Advertisement
	require.NoError(t, json.Unmarshal(body, &resp), "Не удалось распарсить ответ созданного объявления")
	return resp
}

// GetAdvertisementByID получает объявление по ID и проверяет ожидаемый статус-код.
func GetAdvertisementByID(t *testing.T, token, adID string, expectedStatus int) models.Advertisement {
	t.Helper()
	raw := advertisementsClient.HttpGetAdvertisementByID(t, token, adID)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа GET /advertisement")
	require.Equalf(t, expectedStatus, raw.StatusCode, "Некорректный HTTP-статус GET /advertisement. Ответ: %s", string(body))

	var resp models.Advertisement
	require.NoError(t, json.Unmarshal(body, &resp), "Ошибка десериализации объявления")
	return resp
}

// GetAdvertisementPhotos получает список фотографий объявления и проверяет ожидаемый статус-код.
func GetAdvertisementPhotos(t *testing.T, token, adID string, expectedStatus int) []models.Photo {
	t.Helper()
	raw := advertisementsClient.HttpGetAdvertisementPhotos(t, token, adID)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа фотографий")
	require.Equalf(t, expectedStatus, raw.StatusCode, "Некорректный HTTP-статус получения фотографий. Ответ: %s", string(body))

	var photos []models.Photo
	require.NoError(t, json.Unmarshal(body, &photos), "Ошибка десериализации списка фотографий")
	return photos
}

// DeleteAdvertisement удаляет объявление по ID и проверяет ожидаемый статус-код.
func DeleteAdvertisement(t *testing.T, token, adID string, expectedStatus int) {
	t.Helper()
	raw := advertisementsClient.HttpDeleteAdvertisement(t, token, adID)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа при удалении объявления")
	require.Equalf(t, expectedStatus, raw.StatusCode, "Некорректный HTTP-статус при удалении объявления %s. Ответ: %s", adID, string(body))
}

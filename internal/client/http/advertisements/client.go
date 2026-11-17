package advertisements

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// CreateAdvertisementFormParams описывает параметры multipart-формы для создания объявления.
type CreateAdvertisementFormParams struct {
	Title       string
	Description string
	Price       *float64
	Quantity    *int
	PhotoPaths  []string
}

// HttpGetAdvertisements выполняет GET-запрос для поиска объявлений по названию.
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpPostAdvertisement отправляет multipart/form-data запрос на создание нового объявления.
func HttpPostAdvertisement(t *testing.T, token string, params CreateAdvertisementFormParams) *http.Response {
	req := apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", params.Title).
		MultipartFormData("description", params.Description)

	if params.Price != nil {
		req = req.MultipartFormData("price", strconv.FormatFloat(*params.Price, 'f', -1, 64))
	}

	if params.Quantity != nil {
		req = req.MultipartFormData("quantity", strconv.Itoa(*params.Quantity))
	}

	if len(params.PhotoPaths) > 0 {
		req = req.MultipartFile("photos", params.PhotoPaths...)
	}

	return req.Expect(t).End().Response
}

// HttpGetAdvertisementByID запрашивает объявление по его ID через query-параметр.
func HttpGetAdvertisementByID(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPath).
		Query("id", adID).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos запрашивает список фотографий объявления по его ID.
func HttpGetAdvertisementPhotos(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPhotosPath(adID)).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement удаляет объявление по его ID.
func HttpDeleteAdvertisement(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", adID).
		Expect(t).
		End().Response
}

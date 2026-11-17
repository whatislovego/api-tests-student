package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/managers/auth/models"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// Register регистрирует нового пользователя и возвращает данные авторизации с токеном.
func Register(t *testing.T, email, password, firstName string) models.LoginOkResponse {
	t.Helper()
	req := models.RegisterRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
	}

	raw := auth.HttpPostAuthRegister(t, req)
	defer func() { _ = raw.Body.Close() }()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа при регистрации")
	require.Equalf(t, http.StatusCreated, raw.StatusCode, "Некорректный статус при регистрации. Ответ: %s", string(body))

	var resp models.LoginOkResponse
	require.NoError(t, json.Unmarshal(body, &resp), "Ошибка десериализации ответа регистрации")
	return resp
}

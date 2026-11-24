package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// AssertErrorResponse проверяет наличие и точные значения полей error и message в теле ошибки.
func AssertErrorResponse(t *testing.T, body []byte, expectedError, expectedMessage, operation string) {
	t.Helper()

	actualError := gjson.GetBytes(body, "error").String()
	actualMessage := gjson.GetBytes(body, "message").String()

	require.NotEmptyf(t, actualError, "Операция %s: поле error пустое. Ответ сервера: %s", operation, string(body))
	require.Equalf(t, expectedError, actualError, "Операция %s: некорректное поле error. Ответ сервера: %s", operation, string(body))

	require.NotEmptyf(t, actualMessage, "Операция %s: поле message пустое. Ответ сервера: %s", operation, string(body))
	require.Equalf(t, expectedMessage, actualMessage, "Операция %s: некорректное поле message. Ответ сервера: %s", operation, string(body))
}

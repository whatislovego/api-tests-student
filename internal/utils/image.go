package utils

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// minimalPNGBase64 содержит валидный PNG-файл размером 1x1 пиксель.
const minimalPNGBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="

// CreateTempPNGFile создаёт временный валидный PNG-файл в директории теста и возвращает путь к нему.
func CreateTempPNGFile(t *testing.T, fileName string) string {
	t.Helper()

	data, err := base64.StdEncoding.DecodeString(minimalPNGBase64)
	require.NoError(t, err, "Не удалось декодировать тестовый PNG")

	filePath := filepath.Join(t.TempDir(), fileName)
	err = os.WriteFile(filePath, data, 0o600)
	require.NoError(t, err, "Не удалось записать временный PNG-файл")

	return filePath
}

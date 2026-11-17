package assertions

import (
	"testing"

	favModels "api-tests-template/internal/managers/favorites/models"

	"github.com/stretchr/testify/require"
)

// ContainsFavoriteByAdID проверяет, что объявление с указанным adID присутствует в списке избранного.
func ContainsFavoriteByAdID(t *testing.T, items []favModels.FavoriteItem, adID string) {
	t.Helper()
	for _, item := range items {
		if item.Advertisement.ID == adID {
			return
		}
	}
	require.Failf(t, "Объявление не найдено в избранном", "Ожидалось наличие объявления с ID %s в списке избранного", adID)
}

// NotContainsFavoriteByAdID проверяет, что объявление с указанным adID отсутствует в списке избранного.
func NotContainsFavoriteByAdID(t *testing.T, items []favModels.FavoriteItem, adID string) {
	t.Helper()
	for _, item := range items {
		require.NotEqualf(t, adID, item.Advertisement.ID, "Объявление с ID %s не должно присутствовать в избранном", adID)
	}
}

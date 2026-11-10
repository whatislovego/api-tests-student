package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(
	t *testing.T,
	items []models.Advertisement,
	title string,
) {
	found := false
	for _, item := range items {
		if item.Title == title {
			found = true
		}
	}
	require.True(t, found, "Объявление с названием %s не найдено", title)
}

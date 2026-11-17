package path

import "fmt"

const (
	AuthLoginPath        = "/auth/login"
	AuthRegisterPath     = "/auth/register"
	MyAdvertisementsPath = "/my/advertisements"
	AdvertisementsPath   = "/advertisements"
	AdvertisementPath    = "/advertisement"
	FavoritesPath        = "/favorites"
)

// AdvertisementPhotosPath возвращает путь для получения или загрузки фотографий объявления по его ID.
func AdvertisementPhotosPath(adID string) string {
	return fmt.Sprintf("%s/%s/photos", AdvertisementsPath, adID)
}

// FavoriteByAdIDPath возвращает путь для удаления объявления из избранного по его ID.
func FavoriteByAdIDPath(adID string) string {
	return fmt.Sprintf("%s/%s", FavoritesPath, adID)
}

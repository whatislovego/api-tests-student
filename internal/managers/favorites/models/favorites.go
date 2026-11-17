package models

import adModels "api-tests-template/internal/managers/advertisements/models"

type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

type FavoriteItem struct {
	FavoriteID    string                 `json:"favorite_id"`
	AddedAt       string                 `json:"added_at"`
	Advertisement adModels.Advertisement `json:"advertisement"`
}

type FavoriteListResponse struct {
	Items      []FavoriteItem `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}

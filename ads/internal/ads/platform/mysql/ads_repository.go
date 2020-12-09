package mysql

import (
	"database/sql"
	"fmt"

	"github.com/schoren/example-adserver/ads/internal/ads/actions"
	"github.com/schoren/example-adserver/pkg/types"
)

func NewAdsRepository(db *sql.DB) *AdsRepository {
	return &AdsRepository{db}
}

var _ actions.ActiveAdGetter = &AdsRepository{}
var _ actions.CreatePersister = &AdsRepository{}
var _ actions.UpdatePersister = &AdsRepository{}

type AdsRepository struct {
	db *sql.DB
}

func (r *AdsRepository) Create(ad types.Ad) (types.Ad, error) {
	stmt, err := r.db.Prepare("INSERT INTO ads (image_url, clickthrough_url) VALUES (?, ?)")
	if err != nil {
		return types.Ad{}, fmt.Errorf("Error preparing insert ad query: %w", err)
	}

	res, err := stmt.Exec(ad.ImageURL, ad.ClickThroughURL)
	if err != nil {
		return types.Ad{}, fmt.Errorf("Error creating ad: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return types.Ad{}, fmt.Errorf("Error fetching new ad id: %w", err)
	}

	newAd := types.Ad{
		ID:              int(id),
		ImageURL:        ad.ImageURL,
		ClickThroughURL: ad.ClickThroughURL,
	}

	return newAd, nil
}

func (r *AdsRepository) Update(ad types.Ad) error {
	stmt, err := r.db.Prepare("UPDATE ads SET image_url = ?, clickthrough_url = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("Error preparing update ad query: %w", err)
	}

	_, err = stmt.Exec(ad.ImageURL, ad.ClickThroughURL, ad.ID)
	if err != nil {
		return fmt.Errorf("Error updating ad: %w", err)
	}

	return nil
}

func (r *AdsRepository) GetActive() ([]types.Ad, error) {
	results, err := r.db.Query("SELECT id, image_url, clickthrough_url FROM ads")
	if err != nil {
		return []types.Ad{}, fmt.Errorf("Error fetching active ads: %w", err)
	}

	ads := []types.Ad{}
	for results.Next() {
		var ad types.Ad
		err = results.Scan(&ad.ID, &ad.ImageURL, &ad.ClickThroughURL)
		if err != nil {
			return []types.Ad{}, fmt.Errorf("Error reading active ad: %w", err)
		}

		ads = append(ads, ad)
	}

	return ads, nil
}

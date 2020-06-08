package mysql

import (
	"database/sql"
	"fmt"

	"github.com/schoren/example-adserver/types"
)

// AdsRepository is a MySQL repository for ads
type AdsRepository struct {
	db *sql.DB
}

func NewAdsRepository(db *sql.DB) *AdsRepository {
	return &AdsRepository{db}
}

// Create persist an ad into the database and returns a new ad with the corresponding ID
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

// Update persist ad changes into the database
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

// GetActive returns all currently active ads
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

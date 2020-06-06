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
		return fmt.Errorf("Error updateing ad: %w", err)
	}

	return nil
}

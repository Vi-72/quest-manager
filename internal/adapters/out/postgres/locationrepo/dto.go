package locationrepo

import "time"

// LocationDTO is the database model for Location.
type LocationDTO struct {
	ID        string  `gorm:"primaryKey"`
	Latitude  float64 `gorm:"not null;index:idx_location_coords"`
	Longitude float64 `gorm:"not null;index:idx_location_coords"`
	Address   *string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (LocationDTO) TableName() string {
	return "locations"
}

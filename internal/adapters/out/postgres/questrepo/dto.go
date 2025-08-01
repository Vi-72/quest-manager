package questrepo

import "time"

// QuestDTO is the database model for Quest.
type QuestDTO struct {
	ID              string `gorm:"primaryKey"`
	Title           string
	Description     string
	Difficulty      string
	Reward          int // Reward level from 1 to 5
	DurationMinutes int // Store duration in minutes

	// Denormalized coordinates (main data for performance)
	TargetLatitude     float64 `gorm:"index:idx_target_location"`
	TargetLongitude    float64 `gorm:"index:idx_target_location"`
	ExecutionLatitude  float64 `gorm:"index:idx_execution_location"`
	ExecutionLongitude float64 `gorm:"index:idx_execution_location"`

	// Optional references to location directory
	TargetLocationID    *string `gorm:"index"` // FK to quest_locations
	ExecutionLocationID *string `gorm:"index"` // FK to quest_locations

	Equipment string  // stored as comma-separated string
	Skills    string  // stored as comma-separated string
	Status    string  `gorm:"index"`
	Creator   string  `gorm:"index"`
	Assignee  *string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (QuestDTO) TableName() string {
	return "quests"
}

// QuestWithAddressDTO extends QuestDTO for JOIN queries with addresses
type QuestWithAddressDTO struct {
	QuestDTO
	TargetAddress    *string `gorm:"column:target_address"`
	ExecutionAddress *string `gorm:"column:execution_address"`
}

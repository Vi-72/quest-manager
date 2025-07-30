package questrepo

import "time"

// QuestDTO is the database model for Quest.
type QuestDTO struct {
	ID                 string `gorm:"primaryKey"`
	Title              string
	Description        string
	Difficulty         string
	Reward             string
	TargetLatitude     float64 `gorm:"index:idx_target_location"`
	TargetLongitude    float64 `gorm:"index:idx_target_location"`
	ExecutionLatitude  float64 `gorm:"index:idx_execution_location"`
	ExecutionLongitude float64 `gorm:"index:idx_execution_location"`
	Equipment          string  // stored as comma-separated string
	Skills             string  // stored as comma-separated string
	Status             string  `gorm:"index"`
	Creator            string  `gorm:"index"`
	Assignee           *string `gorm:"index"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (QuestDTO) TableName() string {
	return "quests"
}

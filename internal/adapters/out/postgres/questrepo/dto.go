package questrepo

import "time"

// QuestDTO is the database model for Quest.
type QuestDTO struct {
	ID                 string `gorm:"primaryKey"`
	Title              string
	Description        string
	Difficulty         string
	Reward             string
	TargetLatitude     float64
	TargetLongitude    float64
	ExecutionLatitude  float64
	ExecutionLongitude float64
	Equipment          string // stored as comma-separated string
	Skills             string // stored as comma-separated string
	Status             string
	Creator            string
	Assignee           *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (QuestDTO) TableName() string {
	return "quests"
}

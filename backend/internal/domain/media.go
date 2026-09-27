package domain

import "time"

type MediaAsset struct {
	ID                  string     `json:"id"`
	Title               string     `json:"title"`
	Description         string     `json:"description,omitempty"`
	AssetType           string     `json:"asset_type"` // 'video', 'image', 'audio'
	URL                 string     `json:"url"`
	ThumbnailURL        string     `json:"thumbnail_url,omitempty"`
	DistributionType    string     `json:"distribution_type"` // 'public', 'reward', 'promotion', 'advertising'
	RewardRequirementID *string    `json:"reward_requirement_id,omitempty"`
	EventID             *string    `json:"event_id,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}

type MediaRepository interface {
	Create(asset *MediaAsset) error
	GetByID(id string) (*MediaAsset, error)
	List() ([]*MediaAsset, error)
	ListByDistribution(distributionType string) ([]*MediaAsset, error)
	ListByRewardRequirement(requirementID string) ([]*MediaAsset, error)
	ListByEvent(eventID string) ([]*MediaAsset, error)
	Update(asset *MediaAsset) error
	Delete(id string) error
}

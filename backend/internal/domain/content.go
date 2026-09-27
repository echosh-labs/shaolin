package domain

import "time"

type ContentItem struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	ItemType  string    `json:"item_type"` // 'news', 'blog', 'faq', 'resource'
	Category  *string   `json:"category,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PublicEvent struct {
	ID                  string     `json:"id"`
	Title               string     `json:"title"`
	Slug                string     `json:"slug"`
	MainImage           *string    `json:"main_image,omitempty"`
	StartDate           string     `json:"start_date"` // YYYY-MM-DD
	EndDate             string     `json:"end_date"`   // YYYY-MM-DD
	StartTime           *string    `json:"start_time,omitempty"`
	EndTime             *string    `json:"end_time,omitempty"`
	FeeCents            int        `json:"fee_cents"`
	MaxParticipants     int        `json:"max_participants"`
	SpotsAvailable      int        `json:"spots_available"`
	RegistrationRequired bool       `json:"registration_required"`
	RegistrationDeadline *string    `json:"registration_deadline,omitempty"`
	LocationName        *string    `json:"location_name,omitempty"`
	LocationAddress     *string    `json:"location_address,omitempty"`
	LocationDetails     *string    `json:"location_details,omitempty"`
	Description         *string    `json:"description,omitempty"`
}

type PublicEventRegistration struct {
	ID             string    `json:"id"`
	EventID        string    `json:"event_id"`
	UserID         *string   `json:"user_id,omitempty"`
	GuestFirstName *string   `json:"guest_first_name,omitempty"`
	GuestLastName  *string   `json:"guest_last_name,omitempty"`
	GuestEmail     *string   `json:"guest_email,omitempty"`
	GuestPhone     *string   `json:"guest_phone,omitempty"`
	Status         string    `json:"status"` // 'pending', 'confirmed', 'cancelled'
	RegisteredAt   time.Time `json:"registered_at"`
	PaidCents      int       `json:"paid_cents"`
}

type Waiver struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ContentRepository interface {
	CreateContentItem(item *ContentItem) error
	GetContentItemByID(id string) (*ContentItem, error)
	GetContentItemBySlug(slug string) (*ContentItem, error)
	ListContentItems(itemType string) ([]*ContentItem, error)
	UpdateContentItem(item *ContentItem) error
	DeleteContentItem(id string) error

	CreatePublicEvent(event *PublicEvent, instructorIDs []string, hallIDs []int) error
	GetPublicEventByID(id string) (*PublicEvent, error)
	GetPublicEventBySlug(slug string) (*PublicEvent, error)
	ListPublicEvents() ([]*PublicEvent, error)

	CreatePublicEventRegistration(reg *PublicEventRegistration) error
	GetEventRegistrations(eventID string) ([]*PublicEventRegistration, error)

	CreateWaiver(waiver *Waiver) error
	GetActiveWaiver() (*Waiver, error)
	GetWaiverByID(id string) (*Waiver, error)
}

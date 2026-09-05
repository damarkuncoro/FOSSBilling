package domain

import (
	"context"
	"time"
)

type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusSending   CampaignStatus = "sending"
	CampaignStatusCompleted CampaignStatus = "completed"
)

type MassMailCampaign struct {
	ID        int64          `json:"id"`
	AdminID   int64          `json:"admin_id"`
	Subject   string         `json:"subject"`
	Content   string         `json:"content"`
	Status    CampaignStatus `json:"status"`
	SentCount int            `json:"sent_count"`
	CreatedAt time.Time      `json:"created_at"`
	SentAt    *time.Time     `json:"sent_at,omitempty"`
}

type MassMailRepository interface {
	GetByID(ctx context.Context, id int64) (*MassMailCampaign, error)
	List(ctx context.Context, limit, offset int) ([]*MassMailCampaign, int, error)
	Create(ctx context.Context, campaign *MassMailCampaign) error
	Update(ctx context.Context, campaign *MassMailCampaign) error
}

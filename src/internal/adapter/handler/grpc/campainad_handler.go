package grpc

import (
	"context"
	"fmt"

	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
	"github.com/ehsanshah/campaign-services/src/internal/core/port"
	"github.com/google/uuid"

	// مسیر پکیج پروتو را دقیق چک کنید (طبق go_package فایل پروتو)
	pb "github.com/ehsanshah/campaign-services/src/pkg/pb/camp/v1"
)

type Server struct {
	pb.UnimplementedCampaignServiceAdServer // این نام در فایل پروتو شما بود
	svc                                     port.CampaignAdService
}

func NewServer(svc port.CampaignAdService) *Server {
	return &Server{svc: svc}
}

func (s *Server) CreateCampaign(ctx context.Context, req *pb.CreateCampaignRequestAd) (*pb.CreateCampaignResponse, error) {
	var campType domain.CampaignAdType
	detailsMap := make(map[string]interface{})

	// مدیریت oneof طبق نام‌گذاری‌های جدید
	switch d := req.Details.(type) {

	// برای AdDetails (تبلیغاتی)
	case *pb.CreateCampaignRequestAd_AdDetails:
		campType = domain.CampaignTypeAd
		if d.AdDetails != nil {
			detailsMap["platform"] = d.AdDetails.Platform
			detailsMap["daily_budget"] = d.AdDetails.DailyBudget
			detailsMap["target_url"] = d.AdDetails.TargetUrl
			detailsMap["keywords"] = d.AdDetails.Keywords
		}

	// برای EmailDetails (ایمیلی)
	case *pb.CreateCampaignRequestAd_EmailDetails:
		campType = domain.CampaignTypeMTA
		if d.EmailDetails != nil {
			detailsMap["subject"] = d.EmailDetails.Subject
			detailsMap["template_id"] = d.EmailDetails.TemplateId
			detailsMap["sender_name"] = d.EmailDetails.SenderName
			detailsMap["sender_email"] = d.EmailDetails.SenderEmail
			detailsMap["recipient_list_ids"] = d.EmailDetails.RecipientListIds
		}

	default:
		return nil, fmt.Errorf("campaign details are required")
	}

	// ساخت آبجکت دامین CampaignAd
	cmd := &domain.CampaignAd{
		ID:             uuid.New().String(),
		OrganizationID: req.OrganizationId,
		Name:           req.Name,
		Status:         "DRAFT",
		Type:           campType,
		ScheduledAt:    req.ScheduledAt,
		Details:        detailsMap,
	}

	if err := s.svc.Create(ctx, cmd); err != nil {
		return &pb.CreateCampaignResponse{Success: false}, err
	}

	return &pb.CreateCampaignResponse{
		Id:      cmd.ID,
		Success: true,
	}, nil
}

// نکته: متدهای GetCampaign و UpdateStatus هم باید طبق این الگو پیاده‌سازی شوند.

package grpc

import (
	"context"
	"time"

	"github.com/ehsanshah/campaign-services/src/internal/core/domain"
	"github.com/ehsanshah/campaign-services/src/internal/core/port"
	pb "github.com/ehsanshah/campaign-services/src/pkg/pb/camp/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CampaignHandler struct {
	pb.UnimplementedCampaignsMtaServiceServer // دقت کنید: نام اینترفیس در پروتو ICampaignsService است
	service                                   port.ICampaignService
}

func NewCampaignMtaHandler(service port.ICampaignService) *CampaignHandler {
	return &CampaignHandler{
		service: service,
	}
}

// CreateCampaign

func (h *CampaignHandler) CreateCampaign(ctx context.Context, req *pb.CreateCampaignRequest) (*pb.CampaignResponse, error) {
	// تبدیل درخواست Proto به Domain
	domainCamp := &domain.Campaign{
		AccountID: req.AccountId,
		Name:      req.Name,
		EmailIDs:  req.EmailIds,

		// تبدیل گیرندگان
		Recipients: domain.CampaignRecipient{
			ListIDs:    req.Recipients.ListIds,
			SegmentIDs: req.Recipients.SegmentIds,
		},

		// تبدیل تنظیمات
		Options: domain.CampaignOptions{
			DeliveryOptimization: req.Options.DeliveryOptimization,
			TrackOpens:           req.Options.TrackOpens,
			TrackClicks:          req.Options.TrackClicks,
			UseGoogleAnalytics:   req.Options.UseGoogleAnalytics,
			TriggerFrequency:     req.Options.TriggerFrequency,
		},
	}

	created, err := h.service.CreateCampaign(ctx, domainCamp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create campaign: %v", err)
	}

	return &pb.CampaignResponse{Campaign: toProto(created)}, nil
}

// ListCampaigns

func (h *CampaignHandler) ListCampaigns(ctx context.Context, req *pb.ListCampaignsRequest) (*pb.ListCampaignsResponse, error) {
	// نکته: اگر در فایل پروتو limit و offset اضافه کردید، اینجا استفاده کنید
	// فعلاً طبق فایل آپلود شده (که limit نداشت) فقط اکانت را پاس می‌دهیم
	// اما در سرویس limit پیش‌فرض 10 را لحاظ می‌کنیم.
	limit := 10
	offset := 0

	campaigns, err := h.service.ListCampaigns(ctx, req.AccountId, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list campaigns: %v", err)
	}

	var pbCampaigns []*pb.Campaign
	for _, c := range campaigns {
		pbCampaigns = append(pbCampaigns, toProto(c))
	}

	return &pb.ListCampaignsResponse{Campaigns: pbCampaigns}, nil
}

// ---------------------------------------------------------
// Helper: تبدیل Domain به Proto (Mapping پیچیده)
// ---------------------------------------------------------
func toProto(c *domain.Campaign) *pb.Campaign {
	// تبدیل فیلترها (چون آرگومان‌ها جنریک هستند)
	var pbFilters []*pb.FilterCondition
	for _, f := range c.Filters {
		var pbArgs []*structpb.Value
		for _, arg := range f.Args {
			val, _ := structpb.NewValue(arg) // تبدیل هوشمند Go Interface به Proto Value
			pbArgs = append(pbArgs, val)
		}
		pbFilters = append(pbFilters, &pb.FilterCondition{
			Operator: f.Operator,
			Args:     pbArgs,
		})
	}

	// تبدیل ExtraFields
	extraFields, _ := structpb.NewStruct(c.ExtraFields)

	return &pb.Campaign{
		Id:            c.ID,
		AccountId:     c.AccountID,
		Name:          c.Name,
		Status:        c.Status,
		TypeForHumans: c.TypeForHumans,

		// زمان‌بندی‌ها
		CreatedAt:        timeToPb(c.CreatedAt),
		UpdatedAt:        timeToPb(c.UpdatedAt),
		ScheduledFor:     timeToPtrPb(c.ScheduledFor),
		StartedAt:        timeToPtrPb(c.StartedAt),
		FinishedAt:       timeToPtrPb(c.FinishedAt),
		StoppedAt:        timeToPtrPb(c.StoppedAt),
		WinnerSelectedAt: timeToPtrPb(c.WinnerSelectedAt),

		// فلگ‌ها
		IsStopped:             c.IsStopped,
		IsCurrentlySendingOut: c.IsCurrentlySending,
		HasWinner:             c.HasWinner,

		// اشیاء تو در تو
		Recipients: &pb.CampaignRecipient{
			ListIds:      c.Recipients.ListIDs,
			SegmentIds:   c.Recipients.SegmentIDs,
			ListNames:    c.Recipients.ListNames,
			SegmentNames: c.Recipients.SegmentNames,
		},

		Options: &pb.CampaignOptions{
			DeliveryOptimization: c.Options.DeliveryOptimization,
			TrackOpens:           c.Options.TrackOpens,
			TrackClicks:          c.Options.TrackClicks,
			UseGoogleAnalytics:   c.Options.UseGoogleAnalytics,
			TriggerFrequency:     c.Options.TriggerFrequency,
		},

		Stats: &pb.CampaignStats{
			Sent:             c.Stats.Sent,
			OpensCount:       c.Stats.OpensCount,
			UniqueOpensCount: c.Stats.UniqueOpensCount,

			// ✅ اصلاح نهایی:
			// Domain (Value) -> Proto (Float)
			// Domain (Text)  -> Proto (String_)
			OpenRate: &pb.StatsRate{
				Float:   c.Stats.OpenRate.Value,
				String_: c.Stats.OpenRate.Text,
			},

			ClicksCount: c.Stats.ClicksCount,
			ClickRate: &pb.StatsRate{
				Float:   c.Stats.ClickRate.Value,
				String_: c.Stats.ClickRate.Text,
			},

			UnsubscribesCount: c.Stats.UnsubscribesCount,
			UnsubscribeRate: &pb.StatsRate{
				Float:   c.Stats.UnsubscribeRate.Value,
				String_: c.Stats.UnsubscribeRate.Text,
			},

			SpamCount: c.Stats.SpamCount,
			SpamRate: &pb.StatsRate{
				Float:   c.Stats.SpamRate.Value,
				String_: c.Stats.SpamRate.Text,
			},

			HardBouncesCount: c.Stats.HardBouncesCount,
			HardBounceRate: &pb.StatsRate{
				Float:   c.Stats.HardBounceRate.Value,
				String_: c.Stats.HardBounceRate.Text,
			},

			SoftBouncesCount: c.Stats.SoftBouncesCount,
			SoftBounceRate: &pb.StatsRate{
				Float:   c.Stats.SoftBounceRate.Value,
				String_: c.Stats.SoftBounceRate.Text,
			},

			DeliveryRate: c.Stats.DeliveryRate,
		},

		Filters:     pbFilters,
		EmailIds:    c.EmailIDs,
		Warnings:    c.Warnings,
		ExtraFields: extraFields,
	}
}

// توابع کمکی تبدیل زمان
func timeToPb(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func timeToPtrPb(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

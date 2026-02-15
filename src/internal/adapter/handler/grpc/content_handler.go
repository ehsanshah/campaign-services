package grpc

import (
	"context"
	"time"

	"github.com/ehsanshah/campaign-services/src/internal/core/port"
	"github.com/ehsanshah/campaign-services/src/pkg/pb/contents/v1" // مسیر پروتو کانتنت
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type contentGRPCClient struct {
	client contentv1.ContentServiceClient
	conn   *grpc.ClientConn
}

func NewContentGRPCClient(address string) (port.IContentClient, error) {
	// برقراری اتصال با میکروسرویس کانتنت
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &contentGRPCClient{
		client: contentv1.NewContentServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *contentGRPCClient) CreateSnapshot(ctx context.Context, originalID string, campaignID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.CreateSnapshot(ctx, &contentv1.CreateSnapshotRequest{
		OriginalContentId: originalID,
		CampaignId:        campaignID,
	})
	if err != nil {
		return "", err
	}

	return resp.SnapshotId, nil
}

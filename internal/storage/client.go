package storage

// this is for externla object storage R2 (which is s3 comptaiable i,.e we caqn use the same api as aws s3 object storage)

// for more better implementation you can use interfaces

/*


 */

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// client interacts with the object storage as its s3 comptabile api works for both Cloudflare r2 as well as aws s3
// presignedclient for generating presigned url via which user can upload images or download images directly from the object storage short expiry link
// bucket name (bucket is like a folder in which we want to store image or download images from  on the object storage and in bucket they are storage as key value key ientifies theobject be it image , video any asset )
type Client struct {
	client          *s3.Client // pointer so that we dont copy it each time while passing through functions methods
	presignedClient *s3.PresignClient
}

type R2Config struct {
	AccountID    string // AccountID is specific to Cloudflare R2 for aws s dont need it
	AccessKey    string
	AccessSecret string
	Bucket       string
}

func NewR2(ctx context.Context, cfg R2Config) (*Client, error) {

	r2Cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.AccessSecret, "")),
		config.WithRegion("auto"), // Required by SDK but not used by R2
	)
	if err != nil {
		return nil, fmt.Errorf("storage: failed to load storage config: %w", err)
	}

	client := s3.NewFromConfig(r2Cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
	})

	presignClient := s3.NewPresignClient(client)

	return &Client{
		client:          client,
		presignedClient: presignClient,
	}, nil
}

package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageProvider interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error)
	GetPresignedDownloadURL(ctx context.Context, key string, ttl time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
}

type S3StorageConfig struct {
	Endpoint        string
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

type S3StorageProvider struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

func NewS3StorageProvider(cfg S3StorageConfig) (*S3StorageProvider, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("storage bucket is required")
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		if cfg.Endpoint != "" {
			return aws.Endpoint{
				URL:               cfg.Endpoint,
				SigningRegion:     cfg.Region,
				HostnameImmutable: true,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg := aws.Config{
		Region:                      cfg.Region,
		Credentials:                 aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     cfg.AccessKeyID,
				SecretAccessKey: cfg.SecretAccessKey,
			}, nil
		}),
		EndpointResolverWithOptions: customResolver,
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(client)

	return &S3StorageProvider{
		client:        client,
		presignClient: presignClient,
		bucket:        cfg.Bucket,
	}, nil
}

func (s *S3StorageProvider) Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	key = cleanKey(key)
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("s3 upload: %w", err)
	}
	return key, nil
}

func (s *S3StorageProvider) GetPresignedDownloadURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	key = cleanKey(key)
	req, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("presign download url: %w", err)
	}
	return req.URL, nil
}

func (s *S3StorageProvider) Delete(ctx context.Context, key string) error {
	key = cleanKey(key)
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete s3 object: %w", err)
	}
	return nil
}

func cleanKey(k string) string {
	k = strings.TrimSpace(k)
	k = strings.TrimPrefix(k, "/")
	return url.PathEscape(k)
}


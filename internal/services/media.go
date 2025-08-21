package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/sahmaragaev/lunaria-backend/internal/models"
	"github.com/sahmaragaev/lunaria-backend/internal/repositories"
	"github.com/tcolgate/mp3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaService struct {
	s3Client  *s3.Client
	bucket    string
	repo      *repositories.ConversationRepository
	analytics *repositories.AnalyticsRepository
	endpoint  string
}

func NewMediaServiceWithClient(s3Client *s3.Client, bucket string, repo *repositories.ConversationRepository, analytics *repositories.AnalyticsRepository, endpoint string) *MediaService {
	return &MediaService{
		s3Client:  s3Client,
		bucket:    bucket,
		repo:      repo,
		analytics: analytics,
		endpoint:  endpoint,
	}
}

func (m *MediaService) GeneratePresignedUploadURL(ctx context.Context, userID, fileType, format string) (string, string, error) {
	fileID := uuid.New().String()
	var key string
	timestamp := time.Now().UTC()
	year, month, _ := timestamp.Date()
	switch fileType {
	case "photo":
		key = fmt.Sprintf("users/%s/photos/%d/%02d/%s.%s", userID, year, int(month), fileID, format)
	case "voice":
		key = fmt.Sprintf("users/%s/voice/%d/%02d/%s.%s", userID, year, int(month), fileID, format)
	default:
		return "", "", fmt.Errorf("unsupported file type")
	}
	presignClient := s3.NewPresignClient(m.s3Client)
	presignParams := &s3.PutObjectInput{
		Bucket: &m.bucket,
		Key:    &key,
	}
	presigned, err := presignClient.PresignPutObject(ctx, presignParams, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", "", err
	}

	return presigned.URL, fileID, nil
}

func (m *MediaService) ValidateAndExtractMetadata(ctx context.Context, fileID, fileType, format string, data []byte) (*models.MediaMetadata, error) {
	if fileType == "photo" {
		if len(data) > 10*1024*1024 {
			return nil, fmt.Errorf("photo exceeds 10MB size limit")
		}
		mime := mimetype.Detect(data)
		if mime.Is("image/jpeg") || mime.Is("image/png") {
			img, _, err := image.DecodeConfig(bytes.NewReader(data))
			if err != nil {
				return nil, fmt.Errorf("invalid image file")
			}
			thumbBytes, thumbErr := generateThumbnail(data, format)
			var thumbURL *string
			if thumbErr == nil {
				thumbKey := fmt.Sprintf("thumbnails/%s-thumb.%s", fileID, format)
				_, err := m.s3Client.PutObject(ctx, &s3.PutObjectInput{
					Bucket:      &m.bucket,
					Key:         &thumbKey,
					Body:        bytes.NewReader(thumbBytes),
					ContentType: aws.String(mime.String()),
				})
				if err == nil {
					url := fmt.Sprintf("%s/%s/%s", m.endpoint, m.bucket, thumbKey)
					thumbURL = &url
				}
			}
			return &models.MediaMetadata{
				Type:         "photo",
				Format:       format,
				Size:         int64(len(data)),
				Width:        &img.Width,
				Height:       &img.Height,
				ThumbnailURL: thumbURL,
				Status:       "validated",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("unsupported image format")
	} else if fileType == "voice" {
		if len(data) > 10*1024*1024 {
			return nil, fmt.Errorf("voice file exceeds 10MB size limit")
		}
		var duration float64 = 0
		var bitrate int = 0
		if format == "mp3" {
			r := bytes.NewReader(data)
			dec := mp3.NewDecoder(r)
			var f mp3.Frame
			totalsamples := 0
			totalbytes := 0
			for {
				err := dec.Decode(&f, nil)
				if err != nil {
					break
				}
				totalsamples += f.Samples()
				totalbytes += int(f.Size())
			}
			if totalsamples > 0 {
				duration = float64(totalsamples) / 44100.0
			}
			if duration > 0 {
				bitrate = int(float64(totalbytes*8) / duration)
			}
		}
		return &models.MediaMetadata{
			Type:      "voice",
			Format:    format,
			Size:      int64(len(data)),
			Duration:  &duration,
			Bitrate:   &bitrate,
			Status:    "validated",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	return nil, fmt.Errorf("unsupported file type")
}

func generateThumbnail(data []byte, format string) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	thumb := resizeImage(img, 150, 150)
	buf := new(bytes.Buffer)
	if format == "jpeg" || format == "jpg" {
		jpeg.Encode(buf, thumb, nil)
	} else if format == "png" {
		png.Encode(buf, thumb)
	}
	return buf.Bytes(), nil
}

func resizeImage(img image.Image, width, height int) image.Image {
	return resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
}

func (m *MediaService) ModerateContent(ctx context.Context, media *models.MediaMetadata) (bool, error) {
	return true, nil
}

func (m *MediaService) GetMediaMetadataByID(ctx context.Context, id primitive.ObjectID) (*models.MediaMetadata, error) {
	return m.repo.GetMediaMetadataByID(ctx, id)
}

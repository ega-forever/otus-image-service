package interfaces

import (
	"context"
	"github.com/ega-forever/otus-image-service/internal/domain/models"
)

// todo add more actions
type ImageStorage interface {
	SaveImage(ctx context.Context, event *models.Image) (*models.Image, error)
}

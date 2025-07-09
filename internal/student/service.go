package student

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type StudentFinder interface {
	Get(ctx context.Context, id int64) (*Student, error)
}

type Service struct {
	client StudentFinder
	logger *zap.Logger
}

func NewService(c StudentFinder, lg *zap.Logger) *Service {
	return &Service{client: c, logger: lg}
}

func (s *Service) GenerateReport(ctx context.Context, id int64) ([]byte, error) {
	st, err := s.client.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	pdf, err := BuildReport(st)
	if err != nil {
		s.logger.Error("pdf generation failed", zap.Error(err), zap.Int64("studentID", id))
		return nil, fmt.Errorf("failed to build pdf report: %w", err)
	}

	return pdf, nil
}

package service

import "go.uber.org/zap"

type LogService struct {
	logger *zap.Logger
}

func NewLogService(logger *zap.Logger) *LogService {
	return &LogService{
		logger: logger,
	}
}

func (logger LogService) LogError(description string, err error) {
	logger.logger.Error(description, zap.String("error", err.Error()))
}

package service

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (logger LogService) LogError(description string, err error) {

}

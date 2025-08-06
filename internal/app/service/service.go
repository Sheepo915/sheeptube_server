package service

// import "sheeptube/internal/db"

type Service struct {
	VideoService *VideoService
}

// func NewService(queries *db.Queries) *Service {
func NewService() *Service {
	return &Service{
		// VideoService: NewVideoService(queries),
		VideoService: NewVideoService(),
	}
}

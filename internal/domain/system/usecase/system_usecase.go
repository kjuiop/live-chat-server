package usecase

import "live-chat-server/internal/domain/system"

type systemUseCase struct {
	SystemRepository system.Repository
}

func NewSystemUseCase(repository system.Repository) system.UseCase {
	return &systemUseCase{
		SystemRepository: repository,
	}
}

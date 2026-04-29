package part

type service struct {
	repository PartRepository
}

func New(repository PartRepository) *service {
	return &service{
		repository: repository,
	}
}

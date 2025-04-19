package ports

type Repository struct {
	Minio MinioRepository
	Auth  AuthRepository
	User  UserRepository
}

func NewRepository(minioRepo MinioRepository, authRepo AuthRepository, userRepo UserRepository) *Repository {
	return &Repository{
		Minio: minioRepo,
		Auth:  authRepo,
		User:  userRepo,
	}
}

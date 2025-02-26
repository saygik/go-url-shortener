package usecase

func (u *UseCase) GetURL(str string) (string, error) {
	return u.repo.GetURL(str)
}

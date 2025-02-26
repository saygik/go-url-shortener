package usecase

func (u *UseCase) SaveURL(url string, alias string) (interface{}, error) {
	// TODO: put your service call logic here
	// //	return "implement UseCase method GetClients", nil
	return u.repo.SaveURL(url, alias)
}

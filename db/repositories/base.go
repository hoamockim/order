package repositories

type defaultRepository struct {
}

var defaultRepo defaultRepository

func New() *defaultRepository {
	return &defaultRepo
}

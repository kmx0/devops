package repositories

type Repository interface {
	Update(metric, name, value string) error
}

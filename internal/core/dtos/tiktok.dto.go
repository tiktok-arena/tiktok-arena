package dtos

type CreateTiktok struct {
	Name string `validate:"required"`
	URL  string `validate:"required"`
}

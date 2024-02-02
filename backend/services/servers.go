package services

type Server struct {
	ID   int    `validate:"required"`
	name string `validate:"required"`
	host string `validate:"required,ip"`
}

func CreateServer()

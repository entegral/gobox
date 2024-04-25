package main

import (
	"strconv"

	"github.com/entegral/gobox/skeleton"
)

type Sale struct {
	ID            int
	CustomerEmail string
	Items         []string
}

func (s Sale) Pk() (string, string, error) {
	return "pk", "EMAIL#" + s.CustomerEmail, nil
}

func (s Sale) Sk() (string, string, error) {
	return "sk", "ORDER#" + strconv.Itoa(s.ID), nil
}

func (s Sale) KeyFuncs() []skeleton.Generator {
	return []skeleton.Generator{
		s.Pk,
		s.Sk,
	}
}

func main() {
	sale := Sale{
		ID:            123,
		CustomerEmail: "test@gmail.com",
	}
	keys, err := skeleton.DynamoKeyMapV1(sale.KeyFuncs()...)
	if err != nil {
		panic(err)
	}
	for key, value := range keys {
		println(key, value)
	}
}

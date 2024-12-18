package services

import (
	"github.com/tonymj76/mytheresa-test/ent"
	"github.com/tonymj76/mytheresa-test/storage"
)

type RestServiceConfiguration func(rs *RestService) error

type RestService struct {
	DB *ent.Client
}

func NewRestService(cfgs ...RestServiceConfiguration) (*RestService, error) {
	rs := &RestService{}

	for _, cfg := range cfgs {
		if err := cfg(rs); err != nil {
			return nil, err
		}
	}

	return rs, nil
}

func WithCustomDB(client *ent.Client, err error) RestServiceConfiguration {
	return func(rs *RestService) error {
		rs.DB = client
		return err
	}
}

func WithDBSetup() RestServiceConfiguration {
	db, err := storage.NewDB("")
	return WithCustomDB(db, err)
}

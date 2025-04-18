package client

import (
	"fmt"

	"github.com/worldline-go/klient"
)

type Calendar struct {
	klient *klient.Client
}

func New(opts ...klient.OptionClientFn) (*Calendar, error) {
	client, err := klient.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("provider client creation error=%w", err)
	}

	return &Calendar{
		klient: client,
	}, nil
}

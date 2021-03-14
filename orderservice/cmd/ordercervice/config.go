
package main

import "github.com/kelseyhightower/envconfig"

const appID = "orderservice"

type config struct {
	SrvRESTAddress string `envconfig:"serve_rest_address" default:":8000"`
}

func parseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, err
	}
	return c, nil
}


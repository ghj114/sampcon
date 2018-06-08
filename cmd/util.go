package cmd

import (
	"encoding/json"
	"os"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

func setupSpec(context *cli.Context) (*specs.Spec, error) {
	spec, err := loadspec("/home/hyphon/docker/sampleContainer/config.json")
	if err != nil {
		return nil, err
	}
	return spec, nil

}

func loadspec(path string) (spec *specs.Spec, err error) {
	cf, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer cf.Close()
	if err = json.NewDecoder(cf).Decode(&spec); err != nil {
		return nil, err
	}
	return spec, nil
}

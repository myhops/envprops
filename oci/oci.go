package oci

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	// Import the v1 package
)

type response struct {
	Config *Config `json:"config,omitempty"`
}

type Config struct {
	Cmd        []string `json:"Cmd,omitempty"`
	Entrypoint []string `json:"Entrypoint,omitempty"`
}

func FetchConfig(imageRef string) (*Config, error) {
	// Parse the image reference.
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image reference: %w", err)
	}

	// Get the image.
	img, err := remote.Image(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote image: %w", err)
	}

	// Get the image's raw config.
	configJSON, err := img.RawConfigFile() // Corrected line: use img.Config
	if err != nil {
		return nil, fmt.Errorf("failed to get image config: %w", err)
	}

	// Unmarshal to response
	var resp response

	if err := json.Unmarshal(configJSON, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	res := *resp.Config
	return &res, nil
}


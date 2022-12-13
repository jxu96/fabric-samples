package contract

import (
	"encoding/json"
	"fmt"
)

type DatasetMetadata struct {
	ID string `json:"id"`
	// Resource
	Name string `json:"name"`
	Note string `json:"note"`
	// Metadata
	Title                   string   `json:"title"`
	Description             string   `json:"description"`
	ContainsSubnationalData bool     `json:"containsSubnationalData"`
	Source                  string   `json:"source"`
	Organisation            string   `json:"organisation"`
	Maintainer              string   `json:"maintainer"`
	Date                    string   `json:"date"`
	Location                string   `json:"location"`
	FieldNames              []string `json:"fieldNames"`
	FileTypes               []string `json:"fileTypes"`
	NumberOfRows            int      `json:"numberOfRows"`
	License                 string   `json:"license"`
	DefineLicense           string   `json:"defineLicense"`
	Methodology             string   `json:"methodology"`
	DefineMethodology       string   `json:"defineMethodology"`
	UpdateFrequency         string   `json:"updateFrequency"`
	Comments                string   `json:"comments"`
	Tags                    []string `json:"tags"`
	// External access endpoint
	Endpoint string `json:"endpoint"`
}

type DatasetMetadataPublic struct {
	ID string `json:"id"`
	// Resource
	Name string `json:"name"`
	Note string `json:"note"`
	// Metadata
	Title                   string `json:"title"`
	Description             string `json:"description"`
	ContainsSubnationalData bool   `json:"containsSubnationalData"`
	Source                  string `json:"source"`
	Organisation            string `json:"organisation"`
	Maintainer              string `json:"maintainer"`
	Date                    string `json:"date"`
	Location                string `json:"location"`
	NumberOfRows            int    `json:"numberOfRows"`
	License                 string `json:"license"`
	DefineLicense           string `json:"defineLicense"`
	Methodology             string `json:"methodology"`
	DefineMethodology       string `json:"defineMethodology"`
	UpdateFrequency         string `json:"updateFrequency"`
	Comments                string `json:"comments"`
}

type DatasetMetadataInterface interface {
	ToBytes() ([]byte, error)
	FromBytes(bs []byte) error
	Validate() error
}

func (md *DatasetMetadata) ToBytes() ([]byte, error) {
	bs, err := json.Marshal(*md)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode metadata to bytes.\n%v", err)
	}

	return bs, nil
}

func (md *DatasetMetadata) FromBytes(bs []byte) error {
	err := json.Unmarshal(bs, md)
	if err != nil {
		return fmt.Errorf("Failed to decode metadata.\n%v", err)
	}

	return nil
}

func (md *DatasetMetadata) Validate() error {
	return nil
}

func (md *DatasetMetadataPublic) ToBytes() ([]byte, error) {
	bs, err := json.Marshal(*md)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode metadata to bytes.\n%v", err)
	}

	return bs, nil
}

func (md *DatasetMetadataPublic) FromBytes(bs []byte) error {
	err := json.Unmarshal(bs, md)
	if err != nil {
		return fmt.Errorf("Failed to decode metadata.\n%v", err)
	}

	return nil
}

func (md *DatasetMetadataPublic) Validate() error {
	return nil
}

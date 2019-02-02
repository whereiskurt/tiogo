package client

import (
	"encoding/json"
	"github.com/whereiskurt/tiogo/pkg/tenable"
)

// Converter does need any other objects or references
type Converter struct{}

// NewConvert returns a converter, used by the adapter
func NewConvert() (convert Converter) { return }

func (c *Converter) ToVulnExportStatus(raw []byte) (VulnExportStatus, error) {
	var src tenable.VulnExportStatus
	var dst VulnExportStatus

	err := json.Unmarshal(raw, &src)
	if err != nil {
		return dst, err
	}

	dst.Status = src.Status

	for i := range src.Chunks {
		dst.Chunks = append(dst.Chunks, string(src.Chunks[i]))
	}
	for i := range src.ChunksCancelled {
		dst.ChunksCancelled = append(dst.ChunksCancelled, string(src.ChunksCancelled[i]))
	}
	for i := range src.ChunksFailed {
		dst.ChunksFailed = append(dst.ChunksFailed, string(src.ChunksFailed[i]))
	}

	return dst, nil
}

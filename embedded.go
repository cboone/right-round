package rightround

import _ "embed"

//go:embed progress-indicators.json
var progressIndicatorsJSON []byte

// EmbeddedCatalogJSON returns the raw bytes of the embedded progress indicators catalog.
func EmbeddedCatalogJSON() []byte {
	return progressIndicatorsJSON
}

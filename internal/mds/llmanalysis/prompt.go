package llmanalysis

import (
	_ "embed"
	"fmt"

	"github.com/nahyunsama/ceactl/internal/mds/logcompressor"
)

//go:embed prompts/system_prompt.txt
var SystemPrompt string

// BuildUserPrompt renders the device header, compressed group table, and an
// unparsed-repeat-line note into the user-turn prompt sent alongside
// SystemPrompt.
func BuildUserPrompt(device string, result *logcompressor.Result) (string, error) {
	// TODO: device/time-range header + result.WriteReport table + dynamic
	// "N lines were repeat-compressions, totaling ~M occurrences" note
	// derived from result.Unparsed.
	return "", fmt.Errorf("BuildUserPrompt: not implemented")
}

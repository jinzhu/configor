// +build go1.10

package configor

import (
	"encoding/json"
	"io"
	"strings"
)

// unmarshalJSON unmarshals the given data into the config interface.
// If the errorOnUnmatchedKeys boolean is true, an error will be returned if there
// are keys in the data that do not match fields in the config interface.
func unmarshalJSON(data []byte, config interface{}, errorOnUnmatchedKeys bool) error {
	reader := strings.NewReader(string(data))
	decoder := json.NewDecoder(reader)

	if errorOnUnmatchedKeys {
		decoder.DisallowUnknownFields()
	}

	err := decoder.Decode(config)
	if err != nil && err != io.EOF {
		return err
	}
	return nil

}

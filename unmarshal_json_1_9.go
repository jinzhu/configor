// +build !go1.10

package configor

import (
	"encoding/json"
)

// unmarshalJSON unmarshals the given data into the config interface.
// The errorOnUnmatchedKeys boolean is ignored since the json library only has
// support for that feature from go 1.10 onwards.
// If there are any keys in the data that do not match fields in the config
// interface, they will be silently ignored.
func unmarshalJSON(data []byte, config interface{}, errorOnUnmatchedKeys bool) error {

	return json.Unmarshal(data, &config)

}

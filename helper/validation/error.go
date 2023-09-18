package validation

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type (
	// Errors represents the validation errors that are indexed by struct field names, map or slice keys.
	// values are Error or Errors (for map, slice and array error value is Errors).
	ErrorsWithoutKey map[string]error
)

// Error returns the error string of Errors.
func (es ErrorsWithoutKey) Error() string {
	if len(es) == 0 {
		return ""
	}

	keys := make([]string, len(es))
	i := 0
	for key := range es {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	var s strings.Builder
	for i, key := range keys {
		if i > 0 {
			s.WriteString("; ")
		}
		if errs, ok := es[key].(ErrorsWithoutKey); ok {
			_, _ = fmt.Fprintf(&s, "(%v)", errs)
		} else {
			_, _ = fmt.Fprintf(&s, "%v", es[key].Error())
		}
	}
	s.WriteString(".")
	return s.String()
}

// MarshalJSON converts the Errors into a valid JSON.
func (es ErrorsWithoutKey) MarshalJSON() ([]byte, error) {
	errs := map[string]interface{}{}
	for key, err := range es {
		if ms, ok := err.(json.Marshaler); ok {
			errs[key] = ms
		} else {
			errs[key] = err.Error()
		}
	}
	return json.Marshal(errs)
}

// Filter removes all nils from Errors and returns back the updated Errors as an error.
// If the length of Errors becomes 0, it will return nil.
func (es ErrorsWithoutKey) Filter() error {
	for key, value := range es {
		if value == nil {
			delete(es, key)
		}
	}
	if len(es) == 0 {
		return nil
	}
	return es
}

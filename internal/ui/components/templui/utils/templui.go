// templui util templui.go - version: v1.1.0 installed by templui v1.1.0
package utils

import (
	"encoding/json"
	"fmt"
	"maps"
	"time"

	"crypto/rand"

	"github.com/a-h/templ"

	twmerge "github.com/Oudwins/tailwind-merge-go"
)

// TwMerge combines Tailwind classes and resolves conflicts.
// Example: "bg-red-500 hover:bg-blue-500", "bg-green-500" → "hover:bg-blue-500 bg-green-500"
func TwMerge(classes ...string) string {
	return twmerge.Merge(classes...)
}

// TwIf returns value if condition is true, otherwise an empty value of type T.
// Example: true, "bg-red-500" → "bg-red-500"
func If[T comparable](condition bool, value T) T {
	var empty T
	if condition {
		return value
	}
	return empty
}

// TwIfElse returns trueValue if condition is true, otherwise falseValue.
// Example: true, "bg-red-500", "bg-gray-300" → "bg-red-500"
func IfElse[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// MergeAttributes combines multiple Attributes into one.
// Example: MergeAttributes(attr1, attr2) → combined attributes
func MergeAttributes(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		maps.Copy(merged, attr)
	}
	return merged
}

type Attribute struct {
	Key   string
	Value any
}

func DataBind(name string) Attribute {
	return Attribute{
		Key:   "data-bind",
		Value: name,
	}
}

func DataOn(event, action string) Attribute {
	return Attribute{
		Key:   "data-on:" + event,
		Value: action,
	}
}

func DataOnClick(action string) Attribute {
	return DataOn("click", action)
}

func DataMustMarshalSignals(val any) Attribute {
	bs, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}

	return Attribute{
		Key:   "data-signals",
		Value: string(bs),
	}
}

func Attrs(attrs ...Attribute) templ.Attributes {
	res := make(templ.Attributes, len(attrs))
	for _, attr := range attrs {
		res[attr.Key] = attr.Value
	}
	return res
}

func Attr(k string, v any) Attribute {
	return Attribute{Key: k, Value: v}
}

// RandomID generates a random ID string.
// Example: RandomID() → "id-1a2b3c"
func RandomID() string {
	return fmt.Sprintf("id-%s", rand.Text())
}

// ScriptVersion is a timestamp generated at app start for cache busting.
// Used in Script() templates to append ?v=<timestamp> to script URLs.
var ScriptVersion = fmt.Sprintf("%d", time.Now().Unix())

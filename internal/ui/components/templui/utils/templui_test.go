package utils

import (
	"encoding/json"
	"testing"

	"github.com/a-h/templ"

	"github.com/cszczepaniak/cribbly/internal/assert"
)

func TestIfAndIfElse(t *testing.T) {
	assert.Equal(t, "yes", If(true, "yes"))
	assert.Equal(t, "", If(false, "yes"))

	assert.Equal(t, "true", IfElse(true, "true", "false"))
	assert.Equal(t, "false", IfElse(false, "true", "false"))
}

func TestMergeAttributesAndAttrs(t *testing.T) {
	a := Attrs(
		Attr("class", "btn"),
		Attr("data-id", "123"),
	)
	b := Attrs(
		Attr("data-extra", "x"),
	)

	merged := MergeAttributes(a, b)

	assert.Equal(t, templ.Attributes{
		"class":      "btn",
		"data-id":    "123",
		"data-extra": "x",
	}, merged)
}

func TestDataHelpers(t *testing.T) {
	attr := DataBind("name")
	assert.Equal(t, "data-bind", attr.Key)
	assert.Equal(t, "name", attr.Value)

	click := DataOnClick("do-something")
	assert.Equal(t, "data-on:click", click.Key)
	assert.Equal(t, "do-something", click.Value)

	signals := DataMustMarshalSignals(map[string]any{"foo": "bar"})
	assert.Equal(t, "data-signals", signals.Key)

	var decoded map[string]any
	assert.NoError(t, json.Unmarshal([]byte(signals.Value.(string)), &decoded))
	assert.Equal(t, "bar", decoded["foo"])
}


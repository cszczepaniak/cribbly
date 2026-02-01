package components

import (
	"github.com/a-h/templ"

	"github.com/cszczepaniak/cribbly/internal/ui/components/templui/input"
)

type InputProps struct {
	Label        string
	ErrorMessage string
	DataBind     string
	Password     bool
}

func (p InputProps) getInputElemProps() input.Props {
	attrs := make(templ.Attributes)
	if p.DataBind != "" {
		attrs["data-bind"] = p.DataBind
	}
	typ := input.TypeText
	if p.Password {
		typ = input.TypePassword
	}
	return input.Props{
		ID:         "",
		Class:      "",
		Type:       typ,
		Attributes: attrs,
		HasError:   p.ErrorMessage != "",
	}
}

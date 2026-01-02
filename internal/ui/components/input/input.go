package input

type params struct {
	dataBind string
	placeholder string
	typ string
}

type Opt func(params) params

func WithDataBind(name string) Opt {
	return func(p params) params {
		p.dataBind = name
		return p
	}
}

func WithPlaceholder(s string) Opt {
	return func(p params) params {
		p.placeholder = s
		return p
	}
}

func AsPassword() Opt {
	return func(p params) params {
		p.typ = "password"
		return p
	}
}

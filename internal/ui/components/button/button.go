package button

type params struct {
	onClick string
}

type Opt func(params) params

func OnClick(fn string) Opt {
	return func(p params) params {
		p.onClick = fn
		return p
	}
}

package icons

type iconOpt func(iconOpts) iconOpts

func WithStrokeWidth(w float64) iconOpt {
	return func(opts iconOpts) iconOpts {
		opts.strokeWidth = w
		return opts
	}
}

type iconOpts struct {
	strokeWidth float64
}

func resolveOpts(opts ...iconOpt) iconOpts {
	resolved := iconOpts{
		strokeWidth: 1,
	}

	for _, o := range opts {
		resolved = o(resolved)
	}

	return resolved
}

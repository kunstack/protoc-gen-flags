package flags

type Options struct {
	Prefix    []string
	Delimiter string
}

type Option func(*Options)

func WithPrefix(prefix ...string) Option {
	return func(o *Options) {
		o.Prefix = append(o.Prefix, prefix...)
	}
}

func WithDelimiter(delimiter string) Option {
	return func(o *Options) {
		o.Delimiter = delimiter
	}
}

type NameBuilder struct{}

func (n NameBuilder) Build(name string) string {
	return "name"
}

func NewNameBuilder(opts ...Options) NameBuilder {
	return NameBuilder{}
}

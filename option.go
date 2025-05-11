package gostructx

type options struct {
	MaxDepth       int
	ShowTag        bool
	SkipTag        string
	SkipEmpty      bool
	FilterPrefix   string
	MaxSliceMapLen int // slice/map 最大展开数量
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		MaxDepth: 5,
		SkipTag:  "-",
	}
}

// WithMaxDepth 设置最大递归深度
func WithMaxDepth(depth int) Option {
	return func(o *options) {
		o.MaxDepth = depth
	}
}

// WithSkipTag 设置跳过的 tag（比如 json:"-")
func WithSkipTag(tag string) Option {
	return func(o *options) {
		o.SkipTag = tag
	}
}

// WithFilterPrefix 设置只包含特定路径前缀的字段
func WithFilterPrefix(prefix string) Option {
	return func(o *options) {
		o.FilterPrefix = prefix
	}
}

// WithShowTag 设置是否显示 tag
func WithShowTag(show bool) Option {
	return func(o *options) {
		o.ShowTag = show
	}
}

// WithSkipEmpty 设置是否跳过空字段
func WithSkipEmpty(skip bool) Option {
	return func(o *options) {
		o.SkipEmpty = skip
	}
}

// WithMaxSliceMapLen 设置 slice/map 最多展开多少元素
func WithMaxSliceMapLen(max int) Option {
	return func(o *options) {
		o.MaxSliceMapLen = max
	}
}

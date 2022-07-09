// Code generated by "goconfig -type SuccessStatusCode|int,MaxRequestBodyBytes|int64,ResponseContentCharset|string -option -output config_generated.go -configOption Option"; DO NOT EDIT.

package jsonhandler

type ConfigItem[T any] struct {
	modified     bool
	value        T
	defaultValue T
}

func (s *ConfigItem[T]) Set(value T) {
	s.modified = true
	s.value = value
}
func (s *ConfigItem[T]) Get() T {
	if s.modified {
		return s.value
	}
	return s.defaultValue
}
func (s *ConfigItem[T]) Default() T {
	return s.defaultValue
}
func (s *ConfigItem[T]) IsModified() bool {
	return s.modified
}
func NewConfigItem[T any](defaultValue T) *ConfigItem[T] {
	return &ConfigItem[T]{
		defaultValue: defaultValue,
	}
}

type Config struct {
	SuccessStatusCode      *ConfigItem[int]
	MaxRequestBodyBytes    *ConfigItem[int64]
	ResponseContentCharset *ConfigItem[string]
}
type ConfigBuilder struct {
	successStatusCode      int
	maxRequestBodyBytes    int64
	responseContentCharset string
}

func (s *ConfigBuilder) SuccessStatusCode(v int) *ConfigBuilder {
	s.successStatusCode = v
	return s
}
func (s *ConfigBuilder) MaxRequestBodyBytes(v int64) *ConfigBuilder {
	s.maxRequestBodyBytes = v
	return s
}
func (s *ConfigBuilder) ResponseContentCharset(v string) *ConfigBuilder {
	s.responseContentCharset = v
	return s
}
func (s *ConfigBuilder) Build() *Config {
	return &Config{
		SuccessStatusCode:      NewConfigItem(s.successStatusCode),
		MaxRequestBodyBytes:    NewConfigItem(s.maxRequestBodyBytes),
		ResponseContentCharset: NewConfigItem(s.responseContentCharset),
	}
}

func NewConfigBuilder() *ConfigBuilder { return &ConfigBuilder{} }
func (s *Config) Apply(opt ...Option) {
	for _, x := range opt {
		x(s)
	}
}

type Option func(*Config)

func WithSuccessStatusCode(v int) Option {
	return func(c *Config) {
		c.SuccessStatusCode.Set(v)
	}
}
func WithMaxRequestBodyBytes(v int64) Option {
	return func(c *Config) {
		c.MaxRequestBodyBytes.Set(v)
	}
}
func WithResponseContentCharset(v string) Option {
	return func(c *Config) {
		c.ResponseContentCharset.Set(v)
	}
}
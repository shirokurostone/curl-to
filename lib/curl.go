package lib

type KV struct {
	Key   string
	Value string
}

type Form struct {
	Name      string
	Value     string
	TypeValue string
	Filename  string
	Headers   []KV
}

type AuthType int

const (
	AuthNone AuthType = iota
	AuthBasic
	AuthDigest
)

type CurlParam struct {
	URL      string
	Method   string
	Headers  []KV
	Data     []KV
	Form     []Form
	AuthType AuthType
	User     string
	Password string
}

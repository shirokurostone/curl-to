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

type CurlParam struct {
	URL     string
	Method  string
	Headers []KV
	Data    []KV
	Form    []Form
}

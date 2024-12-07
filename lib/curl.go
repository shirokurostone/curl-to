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

type Data struct {
	Type     DataType
	FileName string
	String   string
	Binary   []byte
}

type DataType int

const (
	DataTypeString DataType = iota
	DataTypeFileString
	DataTypeBinary
	DataTypeFileBinary
	DataTypeStdin
)

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
	Data     []Data
	Form     []Form
	AuthType AuthType
	User     string
	Password string
}

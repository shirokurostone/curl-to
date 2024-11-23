package lib

type Form struct {
	Name      string
	Value     string
	TypeValue string
	Filename  string
	Headers   [][2]string
}

type CurlParam struct {
	URL     string
	Method  string
	Headers [][2]string
	Data    [][2]string
	Form    []Form
}

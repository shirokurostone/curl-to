package lib

type CurlParam struct {
	URL     string
	Method  string
	Headers [][2]string
	Data    [][2]string
}

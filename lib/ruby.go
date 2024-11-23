package lib

import (
	"errors"
	"strings"
	"text/template"
)

var errUnknownHttpMethod = errors.New("unknown http method")

func getRequestClassName(method string) (string, error) {
	switch method {
	case "DELETE":
		return "Net::HTTP::Delete", nil
	case "GET":
		return "Net::HTTP::Get", nil
	case "HEAD":
		return "Net::HTTP::Head", nil
	case "PATCH":
		return "Net::HTTP::Patch", nil
	case "POST":
		return "Net::HTTP::Post", nil
	case "PUT":
		return "Net::HTTP::Put", nil
	default:
		return "", errUnknownHttpMethod
	}
}

type templateParams struct {
	URL          string
	RequestClass string
	Headers      [][2]string
	Data         [][2]string
}

func escapeSingleQuoteString(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(value, "\\", "\\\\"), "'", "\\'")
}

func toRubyHash(pairs [][2]string) string {
	sb := new(strings.Builder)
	sb.WriteString("{")
	for i, kv := range pairs {
		sb.WriteString("'")
		sb.WriteString(escapeSingleQuoteString(kv[0]))
		sb.WriteString("' => '")
		sb.WriteString(escapeSingleQuoteString(kv[1]))
		sb.WriteString("'")
		if i != len(pairs)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func GenerateRubyCode(param CurlParam) (string, error) {

	sb := new(strings.Builder)
	tmpl := `require 'net/http'
require 'uri'

url = URI.parse('{{ .URL | escapeSingleQuoteString }}')
req = {{ .RequestClass }}.new(url.request_uri)

{{ range .Headers }}req['{{ index . 0 | escapeSingleQuoteString }}'] = '{{ index . 1 | escapeSingleQuoteString }}'
{{ end }}
{{ if ne .Data nil }}req.set_form_data({{ .Data | toRubyHash }}){{ end }}

http = Net::HTTP.new(url.host, url.port)
http.use_ssl = true if url.is_a?(URI::HTTPS)

res = http.start{ |http|
  http.request(req)
}

puts res.body
`
	t := template.Must(
		template.New("ruby").
			Funcs(template.FuncMap{
				"escapeSingleQuoteString": escapeSingleQuoteString,
				"toRubyHash":              toRubyHash,
			}).
			Parse(tmpl),
	)

	requestClass, err := getRequestClassName(param.Method)
	if err != nil {
		return "", err
	}

	templateParam := templateParams{
		URL:          param.URL,
		RequestClass: requestClass,
		Headers:      param.Headers,
		Data:         param.Data,
	}

	if err := t.Execute(sb, templateParam); err != nil {
		return "", err
	}
	return sb.String(), nil

}

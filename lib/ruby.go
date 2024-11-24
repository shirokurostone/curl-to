package lib

import (
	"errors"
	"fmt"
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

type form struct {
	Name   string
	Value  string
	Params []KV
}

type templateParams struct {
	URL          string
	RequestClass string
	Headers      []KV
	Data         []KV
	DataBinary   string
	Form         []form
	AuthType     AuthType
	User         string
	Password     string
}

func (t templateParams) IsBasic() bool {
	return t.AuthType == AuthBasic
}

func escapeSingleQuoteString(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(value, "\\", "\\\\"), "'", "\\'")
}

func toRubyHash(pairs []KV) string {
	sb := new(strings.Builder)
	sb.WriteString("{")
	for i, kv := range pairs {
		sb.WriteString("'")
		sb.WriteString(escapeSingleQuoteString(kv.Key))
		sb.WriteString("' => '")
		sb.WriteString(escapeSingleQuoteString(kv.Value))
		sb.WriteString("'")
		if i != len(pairs)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func toRubySymbolHash(pairs []KV) string {
	sb := new(strings.Builder)
	sb.WriteString("{")
	for i, kv := range pairs {
		sb.WriteString(":")
		sb.WriteString(kv.Key)
		sb.WriteString(" => '")
		sb.WriteString(escapeSingleQuoteString(kv.Value))
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

{{ range .Headers }}req['{{ .Key | escapeSingleQuoteString }}'] = '{{ .Value | escapeSingleQuoteString }}'
{{ end }}
{{ if .IsBasic }}req.basic_auth('{{ .User | escapeSingleQuoteString }}', '{{ .Password | escapeSingleQuoteString }}'){{end}}
{{ if ne .Data nil }}req.set_form_data({{ .Data | toRubyHash }}){{ end }}
{{ if ne .DataBinary "" }}req.body = '{{ .DataBinary | escapeSingleQuoteString }}'{{ end }}
{{ if ne .Form nil }}req.set_form(
  [
{{ range .Form }}    ['{{ .Name | escapeSingleQuoteString }}', '{{ .Value | escapeSingleQuoteString }}'{{ if ne .Params nil }}, {{ .Params | toRubySymbolHash }}{{ end }}],
{{ end }}  ],
  'multipart/form-data'
)
{{ end }}

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
				"toRubySymbolHash":        toRubySymbolHash,
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
		DataBinary:   param.DataBinary,
		Form:         nil,
		AuthType:     param.AuthType,
		User:         param.User,
		Password:     param.Password,
	}

	if param.AuthType == AuthDigest {
		return "", fmt.Errorf("unsupported auth type: digest")
	}

	if param.Form != nil {
		for _, f := range param.Form {
			fm := form{
				Name:  f.Name,
				Value: f.Value,
			}

			if f.TypeValue != "" {
				fm.Params = append(fm.Params, KV{"content_type", f.TypeValue})
			}
			if f.Filename != "" {
				fm.Params = append(fm.Params, KV{"filename", f.Filename})
			}
			if f.Headers != nil {
				return "", fmt.Errorf("unsupported parameter: headers=")
			}

			templateParam.Form = append(templateParam.Form, fm)
		}
	}

	if err := t.Execute(sb, templateParam); err != nil {
		return "", err
	}
	return sb.String(), nil

}

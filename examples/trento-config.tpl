{
  "node_meta": {
    {{- $first := true }}
    {{- with node }}
    {{- $nodename := .Node.Node }}
    {{- range nodes }}
    {{- if eq .Node $nodename }}
    {{- range ls (print "trento/v0/hosts/" $nodename "/metadata") }}
      {{- if $first }}{{ $first = false }}{{ else }},{{ end }}
      "trento-{{ .Key }}": "{{ .Value }}"
    {{- end }}
    {{- end }}
    {{- end }}
    {{- end }}
  }
}

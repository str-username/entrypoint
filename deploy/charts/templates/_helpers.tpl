{{/* vim: set filetype=mustache: */}}

{{/* Expand the name of the chart. */}}

{{- define "app.name"}}
{{- default .Values.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Create chart name and version as used by the chart label. */}}

{{- define "app.chart" -}}
{{- printf "%s-%s" .Values.name .Chart.Version | replace "+" "_" | trunc 64 | trimSuffix "-" }}
{{- end}}

{{/* Common chart labels. */}}

{{- define "app.labels" -}}
{{ include "app.selectorLabels" . }}
helm.sh/chart: {{ include "app.chart" . }}
app.kubernetes.io/version: {{ .Chart.Version | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/* Common selector labels. */}}

{{- define "app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

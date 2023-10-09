{{/*
Expand the name of the chart.
*/}}
{{- define "standards-insights.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "standards-insights.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "standards-insights.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "standards-insights.labels" -}}
helm.sh/chart: {{ include "standards-insights.chart" . }}
{{ include "standards-insights.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- range $key, $value := .Values.additionalLabels }}
{{ $key }}: {{ $value | quote }}
{{- end -}}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "standards-insights.selectorLabels" -}}
app.kubernetes.io/name: {{ include "standards-insights.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "standards-insights.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "standards-insights.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the pvct to use
*/}}
{{- define "standards-insights.persistentVolumeClaimName" -}}
{{- if .Values.persistentVolumeClaim.create }}
{{- default (include "standards-insights.fullname" .) .Values.persistentVolumeClaim.name }}
{{- else }}
{{- default "default" .Values.persistentVolumeClaim.name }}
{{- end }}
{{- end }}

{{- define "standards-insights.env" -}}
{{- range $key, $value := .Values.env }}
- name: {{ $key | quote }}
  value: {{ $value | quote }}
{{- end -}}
{{- end -}}
{{- define "murmurations.fullname" -}}
{{- $name := default "murmurations" .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "murmurations.index" -}}
  {{- printf "%s-index" (include "murmurations.fullname" .) -}}
{{- end -}}

{{- define "murmurations.validation" -}}
  {{- printf "%s-validation" (include "murmurations.fullname" .) -}}
{{- end -}}

{{- define "murmurations.library" -}}
  {{- printf "%s-library" (include "murmurations.fullname" .) -}}
{{- end -}}

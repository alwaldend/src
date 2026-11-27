apiVersion: fluxcd.controlplane.io/v1
kind: FluxInstance
metadata:
  name: flux
  namespace: {{ .Release.Namespace }}
spec:
  distribution:
    version: {{ required "distribution.version is required" .Values.distribution.version | quote }}
    registry: {{ required "distribution.registry is required" .Values.distribution.registry | quote }}
    variant: {{ required "distribution.variant is required" .Values.distribution.variant | quote }}
  components:
    - source-controller
    - source-watcher
    - kustomize-controller
    - helm-controller
    - notification-controller
    - image-reflector-controller
    - image-automation-controller
  cluster:
    type: kubernetes
    networkPolicy: true

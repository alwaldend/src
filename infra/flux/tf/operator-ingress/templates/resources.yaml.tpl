apiVersion: v1
kind: ConfigMap
metadata:
  name: flux-operator-vault-ca
  namespace: {{ .Release.Namespace }}
data:
  ca.crt: |-
    {{- .Values.vaultCA | nindent 4 }}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Values.serviceAccount }}
  namespace: {{ .Release.Namespace }}
automountServiceAccountToken: false
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: flux-operator-token-request
  namespace: {{ .Release.Namespace }}
rules:
  - apiGroups: [""]
    resources: ["serviceaccounts/token"]
    resourceNames: [{{ .Values.serviceAccount | quote }}]
    verbs: ["create"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: flux-operator-token-request
  namespace: {{ .Release.Namespace }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: flux-operator-token-request
subjects:
  - kind: ServiceAccount
    name: {{ .Values.certManagerName }}
    namespace: {{ .Values.certManagerNamespace }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: flux-operator-token-review
rules:
  - apiGroups: ["authentication.k8s.io"]
    resources: ["tokenreviews"]
    verbs: ["create"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: flux-operator-token-review
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: flux-operator-token-review
subjects:
  - kind: ServiceAccount
    name: {{ .Values.serviceAccount }}
    namespace: {{ .Release.Namespace }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: flux-operator-users
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: view
subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: "flux-ui:src_infra_flux_users"
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: {{ .Values.name }}
  namespace: {{ .Release.Namespace }}
spec:
  vault:
    server: {{ .Values.vaultURL | quote }}
    path: {{ printf "%s/issuer/default/sign/%s" .Values.vaultPKI .Values.vaultRole | quote }}
    caBundle: {{ .Values.vaultCA | b64enc | quote }}
    auth:
      kubernetes:
        mountPath: {{ printf "/v1/auth/%s" .Values.vaultAuth | quote }}
        role: {{ .Values.vaultRole }}
        serviceAccountRef:
          name: {{ .Values.serviceAccount }}
          audiences:
            - https://kubernetes.default.svc.cluster.local

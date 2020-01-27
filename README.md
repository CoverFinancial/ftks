# File To Kubernetes Secrets

read a properties files and create or update kubernetes secrets base on it.


## how to use

```
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: test
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: test
    spec:
      serviceAccountName: vault-user
      automountServiceAccountToken: true

      volumes:
      - name: vault-result
        emptyDir:
          medium: Memory
      - name: config-volume
        configMap:
          name: agent-config
      initContainers:
      - name: vault
        image: vault:1.3.1
        command: ["vault", "agent", "--config=/var/vault/config/agent.config"]
        imagePullPolicy: IfNotPresent
        volumeMounts:
        - name: vault-result
          mountPath: /var/vault/out
        - name: config-volume
          mountPath: /var/vault/config
      - name: fkts
        image: c4po/ftks:v1
        env:
          - name: SECRET_NAME
            value: "test"
          - name: SECRET_FILE
            value: "/var/vault/out/vault.txt"
          - name: SECRET_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
        volumeMounts:
        - name: vault-result
          mountPath: /var/vault/out

      containers:
      - name: sleep
        image: python:3.7
        command: ["/bin/sleep", "3650d"]
        imagePullPolicy: IfNotPresent
        volumeMounts:
        - name: vault-result
          mountPath: /var/vault/out
        envFrom:
        - secretRef:
            name: test

```

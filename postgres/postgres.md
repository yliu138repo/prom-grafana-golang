# Installation

```bash
helm install postgres-normal bitnami/postgresql:latest \
    --namespace default \
    --set persistence.enabled=true \
    --set persistence.size=1Gi \
    --set global.postgresql.postgresqlDatabase=mytestdb \
    --set global.postgresql.postgresqlUsername=dbuser \
    --set global.postgresql.postgresqlPassword=password1234
```
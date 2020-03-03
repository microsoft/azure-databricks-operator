#!/bin/bash

# Find all already running kubectl port-forwards and kill them
ps aux | grep [k]ubectl | awk '{print $2}' | xargs kill > /dev/null 2>&1

echo "-------> Verify all pods up"
echo "-------> Wait for Operator to be ready"
# Ensure there are no Databricks Operator pods not in the Ready state
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl -n azure-databricks-operator-system -lcontrol-plane=controller-manager get pods -o jsonpath="$JSONPATH" 2>&1 | grep -q -v "Ready=False"; do sleep 5;echo "--------> waiting for all operators to be Ready"; kubectl get pods -n azure-databricks-operator-system; done
echo "-------> Wait for MockAPI to be ready"
# Ensure there are no MockAPI pods not in the Ready state
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl -n databricks-mock-api -lapp=databricks-mock-api get pods -o jsonpath="$JSONPATH" 2>&1 | grep -q -v "Ready=fALSE"; do sleep 5;echo "--------> waiting for all mocks to be Ready"; kubectl get pods -n databricks-mock-api; done
# Ensure there are no locust pods not in the Ready state
echo "-------> Wait for locust to be ready"
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl -n locust get pods -lapp=locust-loadtest -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 5;echo "--------> waiting for locust to be available"; kubectl get pods -n locust; done

echo "-------> Open port-forwards"
kubectl port-forward service/prom-azure-databricks-operator-grafana -n default 8080:80 &
kubectl port-forward service/prom-azure-databricks-oper-prometheus -n default 9091:9090 &
kubectl port-forward service/locust-loadtest 8089:8089 9090:9090 -n locust &

echo "Browse to locust webui   -> http://localhost:8089/"
echo "Browse to locust metrics -> http://localhost:9090/"
echo "Browse to Prometheus     -> http://localhost:9091/targets"
echo "Browse to Grafana        -> http://localhost:8080/"
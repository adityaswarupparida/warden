include .env
export

.PHONY: setup setup-external delete cloud-provider cluster-create seal-secrets backup-key restore-key sealed-secrets postgres monitoring upgrade-monitoring deploy rollout

##############################################################
# Setup
##############################################################

setup: cluster-create sealed-secrets restore-key postgres monitoring deploy
	@echo "=======================Cluster is ready==============================="

setup-external: cluster-create sealed-secrets seal-secrets monitoring deploy
	@echo "=======================Cluster is ready==============================="

##############################################################
# Cluster
##############################################################

cluster-create:
	kind create cluster --name warden --config deploy/cluster.yml

delete:
	kind delete cluster --name warden

cloud-provider:
	sudo cloud-provider-kind

##############################################################
# Secrets
##############################################################

seal-secrets:
	kubectl create secret generic warden-secret \
		--dry-run=client \
		'--from-literal=postgresql-url=$(POSTGRESQL_URL)' \
		-o yaml | \
	kubeseal \
		--controller-name=sealed-secrets \
		--controller-namespace=kube-system \
		--format yaml > deploy/k8s/sealed-secret.yml

backup-key:
	kubectl get secret -n kube-system -l sealedsecrets.bitnami.com/sealed-secrets-key -o yaml > sealed-secrets-key.yml
	@echo "Key backed up to sealed-secrets-key.yml — do NOT commit this file"

restore-key:
	kubectl apply -f sealed-secrets-key.yml
	kubectl rollout restart deployment sealed-secrets -n kube-system

##############################################################
# Helm
##############################################################

sealed-secrets:
	helm install sealed-secrets oci://registry-1.docker.io/bitnamicharts/sealed-secrets -n kube-system

postgres:
	helm install postgres oci://registry-1.docker.io/bitnamicharts/postgresql -f deploy/helm/values/postgres.yml

monitoring:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	helm install monitoring prometheus-community/kube-prometheus-stack -f deploy/helm/values/monitoring.yml

upgrade-monitoring:
	helm upgrade monitoring prometheus-community/kube-prometheus-stack -f deploy/helm/values/monitoring.yml

##############################################################
# App
##############################################################

deploy:
	kubectl apply -f deploy/k8s/

rollout:
	kubectl rollout restart deployment warden

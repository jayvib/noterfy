
.PHONY: kustomize-dev
kustomize-dev:
	kubectl kustomize ./dev > ./dev/compiled.yaml

.PHONY: deploy-dev
deploy-dev: kustomize-dev
	kubectl apply -f ./dev/compiled.yaml

.PHONY: teardown-dev
teardown-dev:
	kubectl delete -f ./dev/compiled.yaml
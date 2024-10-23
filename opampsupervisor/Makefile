TAG_NAME := ""
SUPERVISOR_MODULE := github.com/open-telemetry/opentelemetry-collector-contrib/cmd/opampsupervisor@$(TAG_NAME)

src/go.sum:
	go install $(SUPERVISOR_MODULE)
	cp -r $$(go env GOMODCACHE)/$(SUPERVISOR_MODULE) .
	chmod +w opampsupervisor@$(TAG_NAME)
	mv opampsupervisor@$(TAG_NAME) src

getsource: src/go.sum

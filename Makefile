PLUGIN_PATH ?= "${HOME}/.terraform.d/plugins"

install:
	@mkdir -p "${PLUGIN_PATH}"
	@go build -o "${PLUGIN_PATH}/terraform-provider-googlecalendar"

deps:
	@dep ensure -update

dev: install
	@terraform init

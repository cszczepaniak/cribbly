#!/usr/bin/bash

non_templ_go_files=$(find . \( \
		-path ./vendor -o \
		-path ./internal/ui/components/templui \
	\) -prune -o \
	-type f \
	-name "*.go" \
	-not -name "*templ.go" \
	-not -name "*.pb.go" \
	-not -name "*.connect.go" \
	-print)
goimports -w -l -local github.com/cszczepaniak/cribbly $non_templ_go_files

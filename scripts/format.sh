#!/usr/bin/bash

non_templ_go_files=$(find . \
	-name "*.go" \
	-not -name "*templ.go" \
	-not -name "*.pb.go" \
	-not -name "*.connect.go" | grep -v templui)
goimports -w -l -local github.com/cszczepaniak/cribbly $non_templ_go_files

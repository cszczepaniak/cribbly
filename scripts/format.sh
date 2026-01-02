#!/usr/bin/bash

non_templ_go_files=$(find . -name "*.go" -not -name "*templ.go")
goimports -w -l -local github.com/cszczepaniak/cribbly $non_templ_go_files

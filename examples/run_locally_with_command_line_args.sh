#!/bin/bash

../dist/bin/github-commenter \
        -token XXXXXXXXXXXXXXXX \
        -owner cloudposse \
        -repo github-commenter \
        -type pr \
        -number 1 \
        -format "My comment\n{{.}}" \
        -comment "+1 LGTM"

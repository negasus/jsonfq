.SUFFIXES:
.PHONY: help yacc

help: ## show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

yacc: ## Build yacc file with goyacc (golang.org/x/tools)
	@goyacc -o parser/y.go parser/parser.go.y

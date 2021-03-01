SOURCE_MAKE=. ./make/make.sh
SHELL = /bin/bash

test-local:
	@${SOURCE_MAKE} && test-local
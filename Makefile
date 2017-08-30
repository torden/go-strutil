VERSION				:= $(shell git describe --tags --always --dirty="-dev")
DATE				:= $(shell date -u '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS		:= -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'
PLATFORM        	:=$(shell uname -a)
CMD_RM          	:=$(shell which rm)
CMD_CC          	:=$(shell which gcc)
CMD_STRIP       	:=$(shell which strip)
CMD_DIFF        	:=$(shell which diff)
CMD_RM          	:=$(shell which rm)
CMD_BASH        	:=$(shell which bash)
CMD_CP          	:=$(shell which cp)
CMD_AR          	:=$(shell which ar)
CMD_RANLIB      	:=$(shell which ranlib)
CMD_MV          	:=$(shell which mv)
CMD_TAIL        	:=$(shell which tail)
CMD_FIND        	:=$(shell which find)
CMD_LDD         	:=$(shell which ldd)
CMD_MKDIR       	:=$(shell which mkdir)
CMD_TEST        	:=$(shell which test)
CMD_SLEEP       	:=$(shell which sleep)
CMD_SYNC        	:=$(shell which sync)
CMD_LN          	:=$(shell which ln)
CMD_ZIP        		:=$(shell which zip)
CMD_MD5SUM      	:=$(shell which md5sum)
CMD_READELF     	:=$(shell which readelf)
CMD_GDB         	:=$(shell which gdb)
CMD_FILE        	:=$(shell which file)
CMD_ECHO        	:=$(shell which echo)
CMD_NM          	:=$(shell which nm)
CMD_GO				:=$(shell which go)
CMD_GOLINT			:=$(shell which golint)
CMD_GOMETALINTER	:=$(shell which gometalinter)
CMD_GOIMPORTS		:=$(shell which goimport)
CMD_MAKE2HELP		:=$(shell which make2help)
CMD_GLIDE			:=$(shell which glide)
CMD_GOVER			:=$(shell which gover)
CMD_GOVERALLS		:=$(shell which goveralls)

PATH_RACE_REPORT="golang-race.report"
PATH_CONVER_PROFILE="go-strutil.coverprofile"

## Setup Enviroment
setup::
	@$(CMD_ECHO)  -e "\033[1;40;32mSetup Build Enviroment.\033[01;m\x1b[0m"
	@$(CMD_GO) get github.com/Masterminds/glide
	@$(CMD_GO) get github.com/Songmu/make2help/cmd/make2help
	@$(CMD_GO) get github.com/davecgh/go-spew/spew
	@$(CMD_GO) get github.com/k0kubun/pp
	@$(CMD_GO) get github.com/alecthomas/gometalinter
	@$(CMD_GO) get github.com/mattn/goveralls
	@$(CMD_GO) get golang.org/x/tools/cmd/cover
	@$(CMD_GO) get github.com/modocache/gover
	@$(CMD_GO) get github.com/dustin/go-humanize
	@$(CMD_GO) get github.com/golang/lint/golint
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run a LintChecker (Normal)
lint: setup
	@$(CMD_ECHO)  -e "\033[1;40;32mRun a LintChecker (Normal).\033[01;m\x1b[0m"
	@$(CMD_GO) vet $$($(CMD_GLIDE) novendor)
	@for pkg in $$($(CMD_GLIDE) novendor -x); do \
		$(CMD_GOLINT)  -set_exit_status $$pkg || exit $$?; \
	done
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run a LintChecker (Strict)
strictlint: setup
	@$(CMD_ECHO)  -e "\033[1;40;32mRun a LintChecker (Strict).\033[01;m\x1b[0m"
	@$(CMD_GOMETALINTER) $$($(CMD_GLIDE) novendor)
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Go Test with Data Race Detection
test::
	@$(CMD_ECHO)  -e "\033[1;40;32mRun Go Test.\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mYou will get a report of data race detection in $(PATH_RACE_REPORT).pid\033[01;m\x1b[0m"
	@GORACE="log_path=$(PATH_RACE_REPORT)" $(CMD_GO) test -v -test.parallel 4 -race -coverprofile=$(PATH_CONVER_PROFILE)
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Send a report of coverage profile to coveralls.io
coveralls::
	@$(CMD_ECHO)  -e "\033[1;40;32mSend a report of coverage profile to coveralls.io.\033[01;m\x1b[0m"
	@$(CMD_GOVERALLS) -coverprofile=$(PATH_CONVER_PROFILE) -service=travis-ci
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Generate a report about coverage
cover: test
	@$(CMD_ECHO)  -e "\033[1;40;32mGenerate a report about coverage.\033[01;m\x1b[0m"
	@$(CMD_GO) tool cover -func=$(PATH_CONVER_PROFILE)
	@$(CMD_GO) tool cover -func=$(PATH_CONVER_PROFILE)  -o $(PATH_CONVER_PROFILE).html
	@$(CMD_ECHO) -e "\033[1;40;36mGenerated a report file : $(PATH_CONVER_PROFILE).html\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Show Help
help::
	@$(CMD_MAKE2HELP) $(MAKEFILE_LIST)

## Run the cmd/example.go for development
run::
	@$(CMD_GO) run  -tags debug cmd/example.go

## Clean-up
clean::
	@$(CMD_ECHO)  -e "\033[1;40;32mClean-up.\033[01;m\x1b[0m"
	@$(CMD_RM) -rfv *.coverprofile *.swp *.core *.html
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

.PHONY: setup deps updeps lint strictlint help test

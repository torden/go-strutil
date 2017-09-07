PKG_NAME=go-strutil

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
CMD_AWK				:=$(shell which awk)
CMD_SED				:=$(shell which sed)
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

PATH_REPORT=report
PATH_RACE_REPORT=$(PKG_NAME).race.report
PATH_CONVER_PROFILE=$(PKG_NAME).coverprofile
PATH_PROF_CPU=$(PKG_NAME).cpu.prof
PATH_PROF_MEM=$(PKG_NAME).mem.prof
PATH_PROF_BLOCK=$(PKG_NAME).block.prof
PATH_PROF_MUTEX=$(PKG_NAME).mutex.prof

VER_GOLANG=$(shell go version | awk '{print $$3}' | sed -e "s/go//;s/\.//g")
GOLANGV18_OVER=$(shell [ "$(VER_GOLANG)" -ge "180" ] && echo 1 || echo 0)

all: setup lint

check::
ifeq ($(GOLANGV18_OVER),1)
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate a Mutex profile.\033[01;m\x1b[0m"
	go test -tags unittest -test.parallel 4 -bench . -mutexprofile=$(PATH_REPORT)/raw/$(PATH_PROF_MUTEX)
else
	@$(CMD_ECHO) -e "Not Support Mutex Profling"
endif

	
removedep::
	@$(CMD_ECHO)  -e "\033[1;40;32mRemove Deps. Pkgs..\033[01;m\x1b[0m"
	@go clean -i -n github.com/Masterminds/glide
	@go clean -i -n github.com/Songmu/make2help/cmd/make2help
	@go clean -i -n github.com/davecgh/go-spew/spew
	@go clean -i -n github.com/k0kubun/pp
	@go clean -i -n github.com/alecthomas/gometalinter
	@go clean -i -n github.com/mattn/goveralls
	@go clean -i -n golang.org/x/tools/cmd/cover
	@go clean -i -n github.com/modocache/gover
	@go clean -i -n github.com/dustin/go-humanize
	@go clean -i -n github.com/golang/lint/golint
	@go clean -i -n github.com/awalterschulze/gographviz
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Install deps packages
dep::
	@$(CMD_ECHO)  -e "\033[1;40;32mSetup Build Enviroment.\033[01;m\x1b[0m"
	@go get github.com/Masterminds/glide
	@go get github.com/Songmu/make2help/cmd/make2help
	@go get github.com/davecgh/go-spew/spew
	@go get github.com/k0kubun/pp
	@go get github.com/alecthomas/gometalinter
	@go get github.com/mattn/goveralls
	@go get golang.org/x/tools/cmd/cover
	@go get github.com/modocache/gover
	@go get github.com/dustin/go-humanize
	@go get github.com/golang/lint/golint
	@go get github.com/awalterschulze/gographviz
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Setup Enviroment (run after dep)
setup: dep
	@$(CMD_ECHO)  -e "\033[1;40;32mInstall Dep. Packages.\033[01;m\x1b[0m"
	@gometalinter --install
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run a LintChecker (Normal)
lint::
	@$(CMD_ECHO)  -e "\033[1;40;32mRun a LintChecker (Normal).\033[01;m\x1b[0m"
	@go vet $$(glide novendor)
	@for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run a LintChecker (Strict)
strictlint::
	@$(CMD_ECHO)  -e "\033[1;40;32mRun a LintChecker (Strict).\033[01;m\x1b[0m"
	@gometalinter $$(glide novendor)
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Go Test with Data Race Detection
test: clean
	@$(CMD_MKDIR) -p $(PATH_REPORT)/raw/ $(PATH_REPORT)/doc/
	@$(CMD_ECHO)  -e "\033[1;40;32mRun Go Test.\033[01;m\x1b[0m"
	@GORACE="log_path=$(PATH_REPORT)/doc/$(PATH_RACE_REPORT)" go test -tags unittest -v -test.parallel 4 -race -coverprofile=$(PATH_REPORT)/raw/$(PATH_CONVER_PROFILE)
	@$(CMD_ECHO) -e "\033[1;40;36mGenerated a report of data race detection in $(PATH_REPORT)/doc/$(PATH_RACE_REPORT).pid\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Send a report of coverage profile to coveralls.io (run after test)
coveralls: test
	@$(CMD_ECHO)  -e "\033[1;40;32mSend a report of coverage profile to coveralls.io.\033[01;m\x1b[0m"
	@goveralls -coverprofile=$(PATH_REPORT)/raw/$(PATH_CONVER_PROFILE) -service=travis-ci
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Generate a report about coverage (run after test)
cover: test
	@$(CMD_ECHO)  -e "\033[1;40;32mGenerate a report about coverage.\033[01;m\x1b[0m"
	@go tool cover -func=$(PATH_REPORT)/raw/$(PATH_CONVER_PROFILE) -o $(PATH_REPORT)/doc/$(PATH_CONVER_PROFILE).txt
	@go tool cover -html=$(PATH_REPORT)/raw/$(PATH_CONVER_PROFILE) -o $(PATH_REPORT)/doc/$(PATH_CONVER_PROFILE).html
	@$(CMD_ECHO) -e "\033[1;40;36mGenerated a report file : $(PATH_REPORT)/doc/$(PATH_CONVER_PROFILE).html\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Profiling (run after clean)
pprof: clean
	@$(CMD_MKDIR) -p $(PATH_REPORT)/raw/ $(PATH_REPORT)/doc/
	@$(CMD_ECHO)  -e "\033[1;40;32mGenerate profiles.\033[01;m\x1b[0m"
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate a CPU profile.\033[01;m\x1b[0m"
	@go test -tags unittest -test.parallel 4 -bench . -cpuprofile=$(PATH_REPORT)/raw/$(PATH_PROF_CPU)
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate a Memory profile.\033[01;m\x1b[0m"
	@go test -tags unittest -test.parallel 4 -bench . -memprofile=$(PATH_REPORT)/raw/$(PATH_PROF_MEM)
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate a Block profile.\033[01;m\x1b[0m"
	@go test -tags unittest -test.parallel 4 -bench . -blockprofile=$(PATH_REPORT)/raw/$(PATH_PROF_BLOCK)
ifeq ($(GOLANGV18_OVER),1)
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate a Mutex profile.\033[01;m\x1b[0m"
	@go test -tags unittest -test.parallel 4 -bench . -mutexprofile=$(PATH_REPORT)/raw/$(PATH_PROF_MUTEX)
else
	@$(CMD_ECHO)  -e "\033[1;40;33mSkip Generate a Mutex profile. (you go version is old)\033[01;m\x1b[0m"
endif
	@$(CMD_MV) -f *.test $(PATH_REPORT)/raw/
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Generate report of profiling (run after pprof)
report: pprof
	@$(CMD_MKDIR) -p $(PATH_REPORT)/raw/ $(PATH_REPORT)/doc/
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate all report in text format.\033[01;m\x1b[0m"
	@go tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_CPU) > $(PATH_REPORT)/doc/$(PATH_PROF_CPU).txt
	@go tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MEM) > $(PATH_REPORT)/doc/$(PATH_PROF_MEM).txt
	@go tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_BLOCK) > $(PATH_REPORT)/doc/$(PATH_PROF_BLOCK).txt
ifeq ($(GOLANGV18_OVER),1)
	@go tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MUTEX) > $(PATH_REPORT)/doc/$(PATH_PROF_MUTEX).txt
endif
	@$(CMD_ECHO)  -e "\033[1;40;33mGenerate all report in pdf format.\033[01;m\x1b[0m"
	@go tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_CPU) > $(PATH_REPORT)/doc/$(PATH_PROF_CPU).pdf
	@go tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MEM) > $(PATH_REPORT)/doc/$(PATH_PROF_MEM).pdf
	@go tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_BLOCK) > $(PATH_REPORT)/doc/$(PATH_PROF_BLOCK).pdf
ifeq ($(GOLANGV18_OVER),1)
	@go tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MUTEX) > $(PATH_REPORT)/doc/$(PATH_PROF_MUTEX).pdf
endif
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Show Help
help::
	@make2help $(MAKEFILE_LIST)

## Clean-up
clean::
	@$(CMD_ECHO)  -e "\033[1;40;32mClean-up.\033[01;m\x1b[0m"
	@$(CMD_RM) -rfv *.coverprofile *.swp *.core *.html *.prof *.test *.report ./$(PATH_REPORT)/*
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

.PHONY: clean cover coveralls help lint pprof report run setup strictlint test dep

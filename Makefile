BINARY=makeitparty
VERSION=0.2.0
BUILD_TIME=`date +%FT%T%z`

.DEFAULT_GOAL: $(BINARY)

LDFLAGS=-ldflags "-X main.build_time=${BUILD_TIME} -X main.version=${VERSION}"

SOURCEDIR=.

SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} $(SOURCES)

deployable: $(SOURCES)
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ubuntu-${BINARY} $(SOURCES)

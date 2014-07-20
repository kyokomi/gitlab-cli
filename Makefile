.PHONY: all

all:
	gox -osarch="darwin/amd64"
	zip gitlab-cli_darwin_amd64.zip gitlab-cli_darwin_amd64

clean:
	rm gitlab-cli_*


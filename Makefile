.PHONY: all

all:
	gox -osarch="darwin/amd64" -output="gitlab-cli"
	zip gitlab-cli_darwin_amd64.zip gitlab-cli

clean:
	rm gitlab-cli_*


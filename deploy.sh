go build -o ./build/gitlab-cli
git add ./build/gitlab-cli
git commit -am "auto release build"
git push -u origin master
git subtree push --prefix build/ origin release


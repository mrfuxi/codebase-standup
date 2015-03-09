docker run --privileged=true --rm -it -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:cross bash -c "go get .; ./cross-compile.sh"

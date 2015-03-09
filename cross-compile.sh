APP_NAME="codebase-standup"

mkdir -p bin
rm ./bin/$APP_NAME*
for GOOS in darwin linux; do
    for GOARCH in 386 amd64; do
        GOOS=$GOOS GOARCH=$GOARCH go build -v -o ./bin/$APP_NAME-$GOOS-$GOARCH
    done
done

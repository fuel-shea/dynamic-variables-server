dynamic-variables-server
=======================

## Setup
1. Ensure the following are installed:
  * Go (and all it needs)
  * godep
  * MongoDB
  * (TODO)

2. Ensure the project is within the `src` folder of `$GOPATH`

## Configuration and Running

1. Copy the example config files found in the `example-configs` folder into whatever directory the compiled binary will be runnin from (usually `$GOPATH/bin/`) and remove the `.example suffixes`.

2. Build and compile the project
```bash
$ godep go install
```
3. Execute the binary
```bash
$ $GOPATH/bin/dynamic-variables-server
```
   You can also add the APP_ENV environment variable to signal to the app how to config itself
```bash
$ APP_ENV={DEVELOPMENT|TEST|PRODUCTION} $GOPATH/bin/dynamic-variables-server
```

## Testing

(TODO)

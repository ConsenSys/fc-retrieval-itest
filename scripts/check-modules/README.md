# Check modules

Script to parse go.mod and display related branch name from Github.

## Run

### From Makefile

```
make check-modules
```

### From Binary

```
./check-modules/check-modules
```

## Dev

### Run

```
go run scripts/check-modules/check-modules.go
```

### Build

```
go build -o scripts/check-modules/check-modules scripts/check-modules/check-modules.go
```

# Check main== modules

Script to parse go.mod and for usage of main modules.

## Run

### From Makefile

```
make check-main-modules
```

### From Binary

```
./check-main-modules/check-main-modules
```

## Dev

### Run

```
go run scripts/check-main-modules/check-main-modules.go
```

### Build

```
go build -o scripts/check-main-modules/check-main-modules scripts/check-main-modules/check-main-modules.go
```

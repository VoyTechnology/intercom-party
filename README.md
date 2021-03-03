# Intercom Party 🥳

Intercom Party invitation system

---

## Description

This tool is designed to be ran on the CLI. It accepts the customer list as a
standard input, and will return the results on the standard out. The parameters
for running can be set using flags.

---

## Building

### Requirements

- Go 1.16+

### Compiling

To compile a binary version of the application, run

```shell
go build ./cmd/party
```

---

## Running

The one line command to filter customers 100km from the Dublin office (default
flags) is:

```shell
cat customers.txt | party | tee output.txt
```

The defaults can howeverbe changed by the following flags:

| Name        | Default   | Description |
| ----------- | --------- | ----------- |
| `-office`   | `Dublin`  | Office to be used to host the party. Currently only `Dublin` is available, but more can be added in `internal/office/office.go`. In the future this can be replaced by office inventory system or some dynamic API call.
| `-distance` | `100km` | Human readable format for writing the maximum distance. This is in a format of `Xkm`, `XkmYm`, `Ym`, where `X` and `Y` are numbers.

---

## Testing

All tests can be ran by the following command

```shell
go test ./...
```

---

## Architecture

### `cmd/party`

This is where the binary lives. It contains the flag parsing, handing stdin and
stdout. All components are integrated in there. The flow it follows is as follows

- Parse Flags
- Get office
- Parse maximum distance
- Get customers from stdin
- Filter customers if they are within the radius of the office
- Sort customers by ID
- Output the new customer list to the stdout

### `internal/customer`

Customer package contains parsing functionality for the customer file, sorting
functions, filtering the customers and finally writing to file.

### `internal/distance`

The distance package contains human readable distance parsing and calculation
of distance between 2 coordinates.

### `internal/office`

Office contains the office coordinates for a list of available offices. It is
currently only Dublin, but more can be added if needed.

# go-emv-code [![CircleCI][circleci-badge]][circleci] [![codecov][codecov-badge]][codecov] [![GoDoc][godoc-badge]][godoc] [![Go Report Card][goreport-badge]][goreport]

[circleci]: https://circleci.com/gh/mercari/go-emv-code/tree/master
[circleci-badge]: https://circleci.com/gh/mercari/go-emv-code/tree/master.svg?style=svg
[codecov]: https://codecov.io/gh/mercari/go-emv-code
[codecov-badge]: https://codecov.io/gh/mercari/go-emv-code/branch/master/graph/badge.svg
[godoc]: https://godoc.org/go.mercari.io/go-emv-code
[godoc-badge]: https://godoc.org/go.mercari.io/go-emv-code?status.svg
[goreport]: https://goreportcard.com/report/go.mercari.io/go-emv-code
[goreport-badge]: https://goreportcard.com/badge/go.mercari.io/go-emv-code

go-emv-code is a Encoder/Decoder implementation for generate EMV<sup>®</sup><sup>[1](#1)</sup> compliant QR Code<sup>[2](#2)</sup> in Go.

## Usage

See [example](https://godoc.org/go.mercari.io/go-emv-code/mpm/#pkg-examples).

## TODO

* Add Encoder/Decoder implementation for Consumer Presented Mode.

## Contribution

Please read the CLA carefully before submitting your contribution to Mercari.
Under any circumstances, by submitting your contribution, you are deemed to accept and agree to be bound by the terms and conditions of the CLA.

https://www.mercari.com/cla/

### Setup environment & Run tests

* requirements
    * Go version must be at least 1.11 (Modules)

1. Testing in local

```
$ make test
```

2. Testing with [Circle CI CLI](https://circleci.com/docs/2.0/local-jobs/)

```
$ circleci build --job test
```

## License

Copyright 2019 Mercari, Inc.

Licensed under the MIT License.

----

<a name="1">1</a>: EMV<sup>®</sup> is a registered trademark in the U.S. and other countries and an unregistered trademark elsewhere. The EMV trademark is owned by EMVCo, LLC.

<a name="2">2</a>: "QR Code" is a registered trademark of DENSO WAVE

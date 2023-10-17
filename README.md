# go-sourceafis
Golang port of SourceAFIS: A library for human fingerprint recognition

go-sourceafis is a pure Golang port of [SourceAFIS](https://sourceafis.machinezoo.com/),
an algorithm for recognition of human fingerprints.
It can compare two fingerprints 1:1 or search a large database 1:N for matching fingerprint.
It takes fingerprint images on input and produces similarity score on output.
Similarity score is then compared to customizable match threshold.

More on [homepage](https://sourceafis.machinezoo.com/net).

## Status

Unstable but maintained. Some APIs are prone to changes. The WSQ encoder found in **utils/encode/wsq** is a port of NBIS' WSQ decoder algorithm, encoder is yet to be done.

## Getting started

See [homepage](https://sourceafis.machinezoo.com/net).

## Documentation

* [XML doc comments](https://github.com/robertvazan/sourceafis-net/tree/master/SourceAFIS)
* [SourceAFIS overview](https://sourceafis.machinezoo.com/)
* [Algorithm](https://sourceafis.machinezoo.com/algorithm)

## License

Distributed under [Apache License 2.0](https://github.com/jtejido/go-sourceafis/blob/master/LICENSE).
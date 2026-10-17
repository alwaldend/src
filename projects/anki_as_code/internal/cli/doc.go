// Package cli defines the Cobra command tree independently of application implementations.
// Commands consume a private path-based application interface using protobuf
// inputs and outputs. They own flags, syntax, codecs, cancellation and streams;
// executable composition constructs the concrete application outside this package.
package cli

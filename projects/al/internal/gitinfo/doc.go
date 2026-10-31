// Package gitinfo projects Git tags and commits into release changelog metadata.
// The release generator owns bundle loading, history selection, and output;
// this package reads the supplied repository without cloning or writing files.
package gitinfo

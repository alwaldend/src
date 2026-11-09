// Package archive provides bounded, read-only access to standard X export ZIPs.
// It parses literal X source fields directly into canonical archivepb records and
// also provides explicit raw YTD reads. It owns no alternate source record model,
// cleanup graph projection, X authorization or mutation. It never extracts files,
// evaluates JavaScript or follows links. Callers keep the source open while an
// Export is in use and access it sequentially.
package archive

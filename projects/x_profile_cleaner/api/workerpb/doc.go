// Package workerpb declares asynchronous archive-script jobs and typed X-operation plans.
// Trusted status metadata binds results to operation, account, time and source.
// Each operation has one declared payload. Backends own authorization and execution;
// cleanup explanations are a separate optional typed report.
package workerpb

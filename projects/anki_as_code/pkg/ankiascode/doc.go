// Package ankiascode provides configuration and native collection application workflows.
// App exposes path-based Config, Col and ColState operations over the owning
// services. Native operations require worker configuration through New.
// collection.Store lends connections, reserves sources and publishes archives;
// exporters, resource/state stores and the reconciler retain their business rules.
// rebuild.Builder owns copied-package staging, publication and cleanup.
package ankiascode

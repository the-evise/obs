package metrics

/*
   Responsibility:
       - Configure Prometheus registry
       - Expose /metrics endpoint
       - Create base instruments
   Owns:
       - Metric naming
       - Label constraints
   Does NOT:
       - Record request metrics directly
*/

type Metrics interface {
	Inc(name string, labels map[string]string)
}

package prometheus

type PrometheusMetricBuilder interface {
	Namespace(string) PrometheusMetricBuilder
	Name(string) PrometheusMetricBuilder
	Help(string) PrometheusMetricBuilder
	ConstLabels(map[string]string) PrometheusMetricBuilder
	LabelNames([]string) PrometheusMetricBuilder
	//BuildCounterVec() *prometheus.CounterVec
	BuildHistogramVec() (HistogramVec, error)
	BuildGougeVec() (GaugeVec, error)
	//BuildHistogram() (HistogramVec, error)
}

type prometheusMetricBuilder struct {
	namespace   string
	name        string
	help        string
	constLabels map[string]string
	labelNames  []string
	bukets      []float64
}

// Namespace and Name are components of the fully-qualified
// name of the Metric (created by joining these components with
// "_"). Only Name is mandatory, the others merely help structuring the
// name. Note that the fully-qualified name of the metric must be a
// valid Prometheus metric name.
func (p *prometheusMetricBuilder) Namespace(value string) PrometheusMetricBuilder {
	p.namespace = value
	return p
}

// Namespace and Name are components of the fully-qualified
// name of the Metric (created by joining these components with
// "_"). Only Name is mandatory, the others merely help structuring the
// name. Note that the fully-qualified name of the metric must be a
// valid Prometheus metric name.
func (p *prometheusMetricBuilder) Name(value string) PrometheusMetricBuilder {
	p.name = value
	return p
}

// Help provides information about this metric.
//
// Metrics with the same fully-qualified name must have the same Help
// string.
func (p *prometheusMetricBuilder) Help(value string) PrometheusMetricBuilder {
	p.help = value
	return p
}

// ConstLabels are used to attach fixed labels to this metric. Metrics
// with the same fully-qualified name must have the same label names in
// their ConstLabels.
func (p *prometheusMetricBuilder) ConstLabels(value map[string]string) PrometheusMetricBuilder {
	p.constLabels = value
	return p
}

// LabelNames only contain the label names. Their label values are variable
// and therefore not part of the Desc. (They are managed within the Metric.)
func (p *prometheusMetricBuilder) LabelNames(value []string) PrometheusMetricBuilder {
	p.labelNames = value
	return p
}

func (p *prometheusMetricBuilder) Bukets(value []float64) PrometheusMetricBuilder {
	p.bukets = value
	return p
}

func (p *prometheusMetricBuilder) BuildHistogram() (HistogramVec, error) {
	histogramVec, err := NewHistogramVec(
		p.namespace,
		p.name,
		p.help,
		p.bukets,
		p.constLabels,
		p.labelNames,
	)

	return histogramVec, err
}

func (p *prometheusMetricBuilder) BuildGougeVec() (GaugeVec, error) {
	histogramVec, err := NewGaugeVec(
		p.namespace,
		p.name,
		p.help,
		p.bukets,
		p.constLabels,
		p.labelNames,
	)

	return histogramVec, err
}

func (p *prometheusMetricBuilder) BuildHistogramVec() (HistogramVec, error) {
	histogramVec, err := NewHistogramVec(
		p.namespace,
		p.name,
		p.help,
		p.bukets,
		p.constLabels,
		p.labelNames,
	)

	return histogramVec, err
}

func NewBuilder() PrometheusMetricBuilder {
	return &prometheusMetricBuilder{
		bukets: []float64{
			.005, .01, .025, .05, .1, .25, .5, 1, 2, 2.5, 5, 7, 10, 15,
		},
	}
}

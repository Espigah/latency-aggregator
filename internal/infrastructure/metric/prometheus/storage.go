package prometheus

import "sync"

//var histogramStorage = make(map[string]HistogramVec)

var metricStorage = sync.Map{}

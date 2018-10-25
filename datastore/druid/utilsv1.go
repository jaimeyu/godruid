package druid

import (
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models/metrics"
)

// DEPRECATED. Remove this file once V1 has been deprecated

// For postprocessing metrics
type PostProcessorV1 interface {
	Apply(input []metrics.TimeseriesEntryResponseV1) []metrics.TimeseriesEntryResponseV1
}

type NoopPostProcessorV1 struct{}

func (pp NoopPostProcessorV1) Apply(input []metrics.TimeseriesEntryResponseV1) []metrics.TimeseriesEntryResponseV1 {
	logger.Log.Debugf("NoopPostProcessorV1.apply called")
	return input
}

type DropKeysPostprocessorV1 struct {
	keysToDrop []string
	countKeys  map[string][]string
}

func (pp DropKeysPostprocessorV1) Apply(input []metrics.TimeseriesEntryResponseV1) []metrics.TimeseriesEntryResponseV1 {
	logger.Log.Debugf("DropKeysPostprocessorV1.apply called with %v, %v, %v", pp.keysToDrop, pp.countKeys, input)
	if len(pp.keysToDrop) > 0 {
		for _, v := range input {
			for countKey, vals := range pp.countKeys {
				if countVal, ok := v.Result[countKey]; ok {

					if intVal, ok := countVal.(float64); ok && intVal == 0 {
						for _, m := range vals {
							delete(v.Result, m)
						}
					}
				}

			}
			for _, k := range pp.keysToDrop {
				delete(v.Result, k)
			}

		}
	}

	return input
}

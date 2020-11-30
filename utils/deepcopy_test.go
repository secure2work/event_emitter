package utils

import (
	"testing"
)

func Benchmark_Map(b *testing.B) {
	copy := make(map[string]interface{})

	copy["global_m1key"] = "global_m1value"
	copy["global_m2key"] = "global_m2value"
	copy["local_m1key"] = "local_m1value"
	copy["local_m2key"] = "local_m2value"

	for i := 0; i < b.N; i++ {
		_, err := Map(copy)
		if err != nil {
			b.Errorf("%s", err.Error())
		}
	}
}

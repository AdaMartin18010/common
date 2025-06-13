package common_test

import (
	"testing"

	cm "common"
)

func TestPath(t *testing.T) {
	t.Logf("ExecutingCurrentFilePath:%+v", cm.ExecutingCurrentFilePath())
	t.Logf("CompiledExectionFilePath:%+v", cm.CompiledExectionFilePath())
	fp, err := cm.ExecutedCurrentFilePath()
	t.Logf("ExecutedCurrentFilePath:%+v,err:%+v", fp, err)
}

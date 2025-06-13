package component_test

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	cmpt "common/model/component"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReflect(t *testing.T) {
	type sche struct {
	}
	type ScheEx sche
	var sc ScheEx
	ts := reflect.TypeOf(sc)
	fmt.Printf(" testing reflet name: %s,\n string:%s,\n PkgPath:%s,\n Kind: %s\n",
		ts.Name(),
		ts.String(),
		ts.PkgPath(),
		ts.Kind())

	require.Equal(t, "ScheEx", ts.Name())
	require.Equal(t, "navigate/common/model/component_test", ts.PkgPath())
	assert.Equal(t, "component_test.ScheEx", ts.String())

	idname := (cmpt.IdName)(fmt.Sprintf("%s_%X", ts.Name(), rand.Intn(int(^uint(0)>>1))))
	fmt.Printf("idname: %s\n", idname)
	time := time.Now()
	fmt.Printf("time: %v\n", time)
}

func TestComponentBase(t *testing.T) {
	tmp := cmpt.NewCptMetaSt(cmpt.IdName("testId"), cmpt.KindName("testKind"), context.TODO())
	//tmp.RecoverWorker = tmp
	//tmp := cmpt.NewComponentBaseData()
	cpbd := (cmpt.CptRoot)(tmp)
	fmt.Printf("IdName: %s,KindName:%s\n", cpbd.Id(), cpbd.Kind())
	assert.Equal(t, cmpt.IdName("testId"), cpbd.Id())
	assert.Equal(t, cmpt.KindName("testKind"), cpbd.Kind())
	assert.Equal(t, false, cpbd.IsRunning())

	err := cpbd.Start()
	if err != nil {
		fmt.Printf("Start() err is :%#v", err)
	}
	assert.Equal(t, true, cpbd.IsRunning())

	err = cpbd.Stop()
	if err != nil {
		fmt.Printf("Start() err is :%#v", err)
	}
	assert.Equal(t, false, cpbd.IsRunning())

	err = cpbd.Finalize()
	if err != nil {
		fmt.Printf("err is :%#v", err)
	}
	assert.Equal(t, false, cpbd.IsRunning())
}

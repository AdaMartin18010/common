package component_test

import (
	"fmt"
	"testing"

	cmp "common/model/component"
)

func TestComponentsNew(t *testing.T) {
	tmps0 := cmp.NewCpts()
	fmt.Printf("Components len: %d\n", tmps0.Len())
	tmps1 := cmp.NewCpts(nil, nil)
	fmt.Printf("Components len: %d\n", tmps1.Len())
}

func TestComponentsNews(t *testing.T) {
	cmp0 := cmp.NewCptMetaSt()
	cmp1 := cmp.NewCptMetaSt()
	cmp2 := cmp.NewCptMetaSt()
	tmps := cmp.NewCpts(cmp0, cmp1, cmp2)
	fmt.Printf("Components len: %d\n", tmps.Len())
	tmps.AddCpts(cmp0, cmp1, nil)
	fmt.Printf("Components len: %d\n", tmps.Len())
}

func TestComponentsOperater(t *testing.T) {
	cmp0 := cmp.NewCptMetaSt()
	cmp1 := cmp.NewCptMetaSt()
	cmp2 := cmp.NewCptMetaSt()
	tmps := cmp.NewCpts(cmp0, nil, cmp1, nil, cmp2)
	fmt.Printf("Components len: %d\n", tmps.Len())

	printFunc := func(c cmp.Cpt) {
		fmt.Printf("%s,IsRunning:%t\n", c.CmptInfo(), c.IsRunning())
	}

	tmps.Each(printFunc)

	fmt.Println("-----components Start()")
	err := tmps.Start()
	if err != nil {
		fmt.Printf("Start() err is :%#v", err)
	}

	tmps.Each(printFunc)

	fmt.Println("-----components Stop()")
	err = tmps.Stop()
	if err != nil {
		fmt.Printf("Stop() err is :%#v", err)
	}

	tmps.Each(printFunc)

	fmt.Println("-----components RemoveCpts(nil)")
	tmps.RemoveCpts(nil)
	fmt.Printf("Components len: %d\n", tmps.Len())

}

var (
	cmp0 = cmp.NewCptMetaSt(cmp.KindName("component0"))
	cmp1 = cmp.NewCptMetaSt(cmp.KindName("component1"))
	cmp2 = cmp.NewCptMetaSt(cmp.KindName("component2"))
	tmps = cmp.NewCpts(cmp0, nil, cmp1, nil, cmp2)
)

func BenchmarkAdd(b *testing.B) {
	fmt.Printf("Before Add Components len: %d\n", tmps.Len())
	for i := 0; i < b.N; i++ {
		tmps.AddCpts(cmp0, nil, cmp1, nil, cmp2)
	}
	fmt.Printf("After Add Components len: %d\n", tmps.Len())
}

func BenchmarkRemove(b *testing.B) {
	fmt.Printf("Before Add Components len: %d\n", tmps.Len())
	for i := 0; i < 100; i++ {
		tmps.AddCpts(cmp0, nil, cmp1, nil, cmp2)
	}
	fmt.Printf("After Add Components len: %d\n", tmps.Len())
	for i := 0; i < b.N; i++ {
		tmps.RemoveCpts(cmp0, cmp1, nil, cmp2)
	}

	ss := ""
	count := 0
	printFunc := func(c cmp.Cpt) {
		ss += fmt.Sprintf("%s,IsRunning:%t\n", c.CmptInfo(), c.IsRunning())
		count += 1
	}
	fmt.Printf("After Remove Components len: %d\n", tmps.Len())
	tmps.Each(printFunc)
	fmt.Printf("%s\ncount:%d\n", ss, count)
}

package chi 

import(
	"testing"
)
func TestRtree(t *testing.T){
    rt:=NewRtree()
	rt.Insert("cat",1)
	rt.Insert("car",2)
	rt.Insert("dog",3)
	if v,ok:=rt.Get("cat");ok{
		t.Log("cat:",v)
	}
}
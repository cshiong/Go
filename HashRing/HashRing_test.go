package HashRing
import (

	"testing"
   "fmt"

)


var  hashring *HashRing
var  key ="testkey"

func init() {
	memcacheServers := []string{"192.168.0.246",
		"192.168.0.247",
		"192.168.0.249"}
	hashring = NewHashRing(3, 100, memcacheServers)

	for _,v := range hashring.VirtualNodes{
		fmt.Printf("vNode's key: %d pNode's name: %s\n",v.HashKey, v.Node)

	}


}

func TestGetNode(t *testing.T) {
	t.Log("TestGetNode")
	server , _:=hashring.GetNode(key)
	t.Logf("the server is: %s\n",server)
	t.Logf("virtualNode should be sorted after getNode() call")
	// the virturalNodes should be sorted.
	for _, h := range hashring.VirtualNodes {
	  t.Logf("vNode's key: %d pNode's name: %s\n",h.HashKey, h.Node)

	}


}

func TestAddNode(t *testing.T) {
    t.Log("TestAddNode")
	t.Logf("add server:192.168.0.248")
	hashring.AddNode("192.168.0.248",3)

	server , _:=hashring.GetNode(key)
	t.Logf("the server is: %s\n",server)

	for _, h := range hashring.VirtualNodes {
		t.Logf("vNode's key: %d pNode's name: %s\n",h.HashKey, h.Node)

	}

}

func TestDeleteNode(t * testing.T){
	t.Log("TestDeleteNode")
    t.Logf("delete server: 192.168.0.246")
	hashring.DeleteNode("192.168.0.246")
	server , _:=hashring.GetNode(key)

	t.Logf("the server is: %s\n",server)
	for _, h := range hashring.VirtualNodes {
		t.Logf("vNode's key: %d pNode's name: %s\n",h.HashKey, h.Node)

	}

}

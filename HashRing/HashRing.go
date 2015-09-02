// HashRing
package HashRing

import (
	"crypto/md5"


	"strconv"

	"sort"
	"encoding/gob"
	"bytes"
	"sync"
	"fmt"
)

type HashKey uint64

//virturalNode using the hasKey as its identity and Node is the physical node this virtual node belong to
type VirtualNode struct{
	HashKey
	Node string
}


func (v VirtualNode) compare (c VirtualNode) int{
	return int(v.HashKey - c.HashKey)
}

//sorted slice based on the HashKey
type VirtualNodes []VirtualNode

// Len is part of sort.Interface.  not using pointer since slice use reference
func (s VirtualNodes) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s VirtualNodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface. i, j is the index of the virtual node in the collection
func (s VirtualNodes) Less(i, j int) bool {
	return s[i].HashKey <s[j].HashKey
}




type HashRing struct {
	sync.RWMutex
	VirtualNodes  //virtual node key mapping to physical nodes,
    Nodes map[string][]VirtualNode //  the physical node have all the virtualNode index, for quick operations based on the physical node
    Capacity int
	Replicates int

}



func NewHashRing(replicates,capacity int, pNodes []string, ) *HashRing{
	 hashRing := &HashRing{
		 VirtualNodes: make([]VirtualNode,0,capacity),
		 Nodes: make(map[string][]VirtualNode), //
		 Replicates: replicates,
		 Capacity:capacity,
	}

    for _,node := range pNodes{
		hashRing.AddNode(node, replicates)
    }
   // sort.Sort(hashRing.VirtualNodes)
	return hashRing
}

func toByteArray(index int, str string) []byte{
	return []byte(strconv.Itoa(index) + "_" + str)
}



func hashVal(key []byte) HashKey {  //to uint64
	bKey :=hashDigest(key)

	return ((HashKey(bKey[3]) << 24) |
	(HashKey(bKey[2]) << 16) |
	(HashKey(bKey[1]) << 8) |
	(HashKey(bKey[0])))
}




func hashDigest(key []byte) []byte {
	m := md5.New()
	m.Write(key)
	return m.Sum(nil)
}


//add all the virtual nodes of this physical nodes to  the ring
func (hashRing * HashRing) AddNode(name string,replicate int){
	for i:=0; i< replicate; i++ {
		bytes := toByteArray(i, name)
		hashK := hashVal(bytes)
		vn :=VirtualNode{hashK,name}
		hashRing.VirtualNodes = append(hashRing.VirtualNodes,vn)
        hashRing.Nodes[name] = append(hashRing.Nodes[name],vn)
	}
}


func (hashRing * HashRing) DeleteNode(name string){
	vs :=hashRing.VirtualNodes
	hashRing.Lock()
	if !sort.IsSorted(vs){
		sort.Sort(vs)
	}
	for _, h := range hashRing.Nodes[name] {
		fmt.Printf("haskKey to remove : %d\n",h.HashKey)
		i := vs.searchPosition(h,0,len(vs)-1)
		hashRing.VirtualNodes = append(hashRing.VirtualNodes[:i],hashRing.VirtualNodes[i+1:]...)
	}
	delete(hashRing.Nodes,name)
	hashRing.Unlock()
}

// find where the key should go which virtual node.(hashKey) and physical node too.
func (h * HashRing ) GetNode( key interface{}) (string,error ){

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return "", err
	}

	hashV := hashVal(buf.Bytes())
	//fmt.Printf("key is: %s HashValue is: %v\n",key, hashV)
    v :=VirtualNode{hashV,""}
	vs :=h.VirtualNodes
	//if !sort.IsSorted(vs){}
	h.Lock()
	if !sort.IsSorted(vs){
		sort.Sort(vs)
	}
    i := vs.searchPosition(v,0,len(vs)-1)
	name :=vs[i].Node
	h.Unlock()
	return name,nil

}




func (hr  VirtualNodes) searchPosition(h VirtualNode,start,end int) int{


	if end - start ==1  {
		return end
	}
	i := start + (end - start)/2
	less := hr[i].compare(h)
	if  less < 0 {
		return hr.searchPosition(h,i,end)
	}else if less > 0{
		return hr.searchPosition(h,start,i-1)
	}else {
		return i
	}
}


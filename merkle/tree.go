package merkle

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
)

// NodeLayer represents a tree layer of hash results
type NodeLayer [][sha1.Size]byte

func (nodes NodeLayer) next() NodeLayer {
	nextLayerSize := len(nodes) / 2
	// if last layer had an odd number of nodes, you need one more than half the previous length
	if len(nodes)%2 != 0 {
		nextLayerSize += 1
	}
	nextLayer := make(NodeLayer, nextLayerSize)
	for i, hash := range nodes {
		h := hash[:]
		// do nothing if first node
		if i == 0 {
			continue
		}
		nextLayerIndex := i / 2

		// the second of each pair is not even due to zero indexing
		secondOfPair := i%2 != 0
		if secondOfPair {
			prevIndex := i - 1
			prevHash := nodes[prevIndex]
			both := append(h, prevHash[:]...)
			nextLayer[nextLayerIndex] = sha1.Sum(both)
		} else if i == len(nodes)-1 {
			// if last node and there were an uneven number, hash against self
			c := make([]byte, sha1.Size)
			copy(c, h)
			both := append(c, h...)
			nextLayer[nextLayerIndex] = sha1.Sum(both)
		}
	}

	return nextLayer
}

func (nodes NodeLayer) Root() [sha1.Size]byte {
	if len(nodes) == 1 {
		return nodes[0]
	}
	layer := nodes
	for len(layer) > 1 {
		layer = layer.next()
	}
	if len(layer) != 1 {
		log.Fatalf("expected a root NodeLayer with one node, got %d", len(layer))
	}
	return layer[0]
}

func From(filenames ...string) (NodeLayer, error) {
	nodes := make(NodeLayer, len(filenames))
	for i, filename := range filenames {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", filename, err)
		}
		nodes[i] = sha1.Sum(bytes)
	}
	return nodes, nil
}

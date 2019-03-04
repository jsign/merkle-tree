package merkle

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidatePartialTree(t *testing.T) {
	req := require.New(t)

	leafIndices := []uint64{3}
	leaves := [][]byte{NewNodeFromUint64(3)}
	proof := [][]byte{
		NewNodeFromUint64(0),
		NewNodeFromUint64(0),
		NewNodeFromUint64(0),
	}
	root, _ := NewNodeFromHex("62b525ec807e21a1fd12d06905d85c4b7bc1feacfa57789d95702f6b69ce129f")
	valid, err := ValidatePartialTree(leafIndices, leaves, proof, root, GetSha256Parent)
	req.NoError(err)
	req.True(valid, "Proof should be valid, but isn't")
}

func TestValidatePartialTreeForRealz(t *testing.T) {
	req := require.New(t)

	leafIndices := []uint64{4}
	leaves := [][]byte{NewNodeFromUint64(4)}
	tree := NewProvingTree(GetSha256Parent, leafIndices)
	for i := uint64(0); i < 8; i++ {
		err := tree.AddLeaf(NewNodeFromUint64(i))
		req.NoError(err)
	}
	root, err := tree.Root() // 4a2ca61d1fd537170785a8575d424634713c82e7392e67795a807653e498cfd0
	req.NoError(err)
	proof, err := tree.Proof() // 05 6b 13
	req.NoError(err)

	valid, err := ValidatePartialTree(leafIndices, leaves, proof, root, GetSha256Parent)
	req.NoError(err)
	req.True(valid, "Proof should be valid, but isn't")

	/***********************************
	|                4a                |
	|       .13.             6c        |
	|    9d      fe      3d     .6b.   |
	|  00  01  02  03 =04=.05. 06  07  |
	***********************************/
}

func TestValidatePartialTreeMulti(t *testing.T) {
	req := require.New(t)

	leafIndices := []uint64{1, 4}
	leaves := [][]byte{
		NewNodeFromUint64(1),
		NewNodeFromUint64(4),
	}
	tree := NewProvingTree(GetSha256Parent, leafIndices)
	for i := uint64(0); i < 8; i++ {
		err := tree.AddLeaf(NewNodeFromUint64(i))
		req.NoError(err)
	}
	root, err := tree.Root() // 4a2ca61d1fd537170785a8575d424634713c82e7392e67795a807653e498cfd0
	req.NoError(err)
	proof, err := tree.Proof() // 05 6b 13
	req.NoError(err)

	valid, err := ValidatePartialTree(leafIndices, leaves, proof, root, GetSha256Parent)
	req.NoError(err)
	req.True(valid, "Proof should be valid, but isn't")

	/***********************************
	|                4a                |
	|        13              6c        |
	|    9d     .fe.     3d     .6b.   |
	| .00.=01= 02  03 =04=.05. 06  07  |
	***********************************/
}

func TestValidatePartialTreeMulti2(t *testing.T) {
	req := require.New(t)

	leafIndices := []uint64{0, 1, 4}
	leaves := [][]byte{
		NewNodeFromUint64(0),
		NewNodeFromUint64(1),
		NewNodeFromUint64(4),
	}
	tree := NewProvingTree(GetSha256Parent, leafIndices)
	for i := uint64(0); i < 8; i++ {
		err := tree.AddLeaf(NewNodeFromUint64(i))
		req.NoError(err)
	}
	root, err := tree.Root() // 4a2ca61d1fd537170785a8575d424634713c82e7392e67795a807653e498cfd0
	req.NoError(err)
	proof, err := tree.Proof() // 05 6b 13
	req.NoError(err)

	valid, err := ValidatePartialTree(leafIndices, leaves, proof, root, GetSha256Parent)
	req.NoError(err)
	req.True(valid, "Proof should be valid, but isn't")

	/***********************************
	|                4a                |
	|        13              6c        |
	|    9d     .fe.     3d     .6b.   |
	| =00==01= 02  03 =04=.05. 06  07  |
	***********************************/
}

func BenchmarkValidatePartialTree(b *testing.B) {
	req := require.New(b)

	leafIndices := []uint64{100, 1000, 10000, 100000, 1000000, 2000000, 4000000, 8000000}
	var leaves [][]byte
	for _, i := range leafIndices {
		leaves = append(leaves, NewNodeFromUint64(i))
	}
	tree := NewProvingTree(GetSha256Parent, leafIndices)
	for i := uint64(0); i < 1<<23; i++ {
		err := tree.AddLeaf(NewNodeFromUint64(i))
		req.NoError(err)
	}
	root, err := tree.Root()
	req.NoError(err)
	proof, err := tree.Proof()
	req.NoError(err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		valid, err := ValidatePartialTree(leafIndices, leaves, proof, root, GetSha256Parent)
		req.NoError(err)
		req.True(valid, "Proof should be valid, but isn't")
	}

	/***********************************
	|                4a                |
	|        13              6c        |
	|    9d     .fe.     3d     .6b.   |
	| =00==01= 02  03 =04=.05. 06  07  |
	***********************************/
}
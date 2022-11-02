package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

type block struct {
	Arr       [10]MerkelTree
	M_index   int
	id        int
	nonce     string
	Prev_hash string
	Curr_hash string
}

type blockchain struct {
	list []*block
}

type M_Node struct {
	Left  *M_Node
	Right *M_Node
	Trx   string
}

type MerkelTree struct {
	root *M_Node
}

func add_node(root *M_Node, data *M_Node) {

	if len(data.Trx) <= len(root.Trx) {
		if root.Left == nil {
			root.Left = data
			return
		} else {
			add_node(root.Left, data)
			return
		}
	} else if len(data.Trx) > len(root.Trx) {
		if root.Right == nil {
			root.Right = data
			return
		} else {
			add_node(root.Right, data)
			return
		}
	}
}

// Adds a block to blockchain.
func addNewBlock(trx string, Arr *[10]MerkelTree, M_index *int) {
	var i int
	i = *M_index
	var Tree *M_Node = Arr[i].root
	var tmp M_Node

	//Adding transaction to the node.
	tmp.Trx = trx
	tmp.Right = nil
	tmp.Left = nil

	if Tree == nil {
		Arr[i].root = &tmp
		i++
		fmt.Println(trx, "inserted")
	} else {
		add_node(Tree, &tmp)
		fmt.Println(trx, "inserted")
		i++
	}
	*M_index = i
}

func Display_node(node *M_Node) {

	if node == nil {
		return
	}
	fmt.Print(node.Trx, ", ")

	if node.Left != nil {
		Display_node(node.Left)
	}

	if node.Right != nil {
		Display_node(node.Right)
	}
}

func UpdateNode(root *M_Node, Prev string, current string) {

	if root == nil {
		return
	}
	if Prev == root.Trx {
		root.Trx = current
	}

	if root.Left != nil {
		UpdateNode(root.Left, Prev, current)
	}

	if root.Right != nil {
		UpdateNode(root.Right, Prev, current)
	}

}

func UpdateAll(Arr *[10]MerkelTree, M_index *int, Prev string, current string) {
	var T_index int
	T_index = *M_index
	var i = 0
	for i = 0; i < T_index+1; i++ {
		var root *M_Node = Arr[i].root
		UpdateNode(root, Prev, current)
	}
	*M_index = T_index
}

func TraverseBlock(root *M_Node, T_Trx *string) {

	if root == nil {
		return
	}
	*T_Trx += root.Trx

	if root.Left != nil {
		TraverseBlock(root.Left, T_Trx)
	}

	if root.Right != nil {
		TraverseBlock(root.Right, T_Trx)
	}

}

func TraverseTillindex(Arr *[10]MerkelTree, M_index *int) string {

	var i = 0
	var T_index int
	T_index = *M_index
	var Trxs string
	for i = 0; i < T_index+1; i++ {
		var root *M_Node = Arr[i].root
		TraverseBlock(root, &Trxs)
	}
	*M_index = T_index
	return Trxs

}

func createMerkelTree(M_tree *[10]MerkelTree, M_index *int) {

	for i := 0; i < 10; i++ {
		M_tree[i].root = nil
	}
	addNewBlock("trx1", M_tree, M_index)
	addNewBlock("trx2", M_tree, M_index)
	addNewBlock("trx3", M_tree, M_index)
	addNewBlock("trx4", M_tree, M_index)
	addNewBlock("trx5", M_tree, M_index)
}

// (1) creates new block. Nonce:=0.
func NewBlock(data int) *block {
	tmp := new(block)
	tmp.id = data
	tmp.nonce = "0" // nonce Initially 0
	tmp.M_index = 0 // root
	createMerkelTree(&tmp.Arr, &tmp.M_index)
	return tmp
}

// (2) A method to find the nonce value for the block.
func Mineblock(chain *blockchain) {
	for i := 0; i < len(chain.list); i++ {
		println("hash:", chain.list[i].Curr_hash, "\n")
		for j := 0; ; j++ {
			Nonce := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.Itoa(j))))
			NValue := Nonce[:3]
			//fmt.Println("nounce:", NValue)
			//fmt.Println(strings.Contains(chain.list[i].Curr_hash, NValue))

			if strings.Contains(chain.list[i].Curr_hash, NValue) == true {
				chain.list[i].nonce = NValue
				break
			}
		} //End for 1
	} // End For-2
}

// (3) Print all the blocks, showing block data such
// as nonce, previous hash, current block hash.
func DisplayBlocks(blocklist *blockchain) {
	for i := 0; i < len(blocklist.list); i++ {
		fmt.Println("---------------------------------------------------")
		fmt.Println("\n\tBlock details:")
		fmt.Printf("Id:%d\n", blocklist.list[i].id)
		fmt.Println("\tMerkel Tree: ")
		DisplayMerkelTree(&blocklist.list[i].Arr, &blocklist.list[i].M_index)
		fmt.Println("nonce: ", blocklist.list[i].nonce)
		fmt.Println("Current block hash: ", blocklist.list[i].Curr_hash)
		fmt.Println("Previous hash: ", blocklist.list[i].Prev_hash)
	}
	fmt.Println("") //New Line
}

// (4) Print all the transactions, showing the transactions and hashes.
func DisplayMerkelTree(Arr *[10]MerkelTree, tree_index *int) {
	var M_index int
	M_index = *tree_index
	var i = 0

	fmt.Printf("<<<<<<<<<<<<<<<<<<<<Tree>>>>>>>>>>>>>>>>>>>>>>\n")
	//fmt.Printf("\n")
	for i = 0; i < M_index+1; i++ {
		var M_tree *M_Node = Arr[i].root
		if M_tree == nil { // Tree is empty
			fmt.Println(i, ": empty")
		} else { // Tree not empty.
			fmt.Print(" ", i, ": ")
			Display_node(M_tree)
		}
	} //End for
	*tree_index = M_index
}

// (5) Function to change one or multiple transactions of the given block ref.
func ChangeBlock(chain *blockchain, block_id int) { // updating on basis of id value as identifier
	flag := false
	for i := 0; i < len(chain.list); i++ {
		if block_id == chain.list[i].id {
			PreviousTrx := "Trx1"
			CurrentTrx := "TrxNew"
			fmt.Println("UPDATED!")
			UpdateAll(&chain.list[i].Arr, &chain.list[i].M_index, PreviousTrx, CurrentTrx)
			flag = true
		}
	}
	if flag == false {
		fmt.Println("ERROR! 404 Block Not Found. :o")
	}
	return
}

// (6) Function to verify blockchain in case any changes are made. The verification will consider the
// changes to the transactions stored in the Merkel tree.
func VerifyChain(chain *blockchain) bool {
	var tmp = ""
	var flag = true
	for i := 0; i < len(chain.list); i++ {
		T_Trx := TraverseTillindex(&chain.list[i].Arr, &chain.list[i].M_index)

		var attr string
		attr += strconv.Itoa(chain.list[i].id)
		attr += T_Trx + chain.list[i].Prev_hash
		sum := sha256.Sum256([]byte(attr))
		tmp = fmt.Sprintf("%x", sum)

		if tmp != chain.list[i].Curr_hash {
			flag = false
			fmt.Printf("Prev Block is changed. With Block id %d", i)
			fmt.Println("")
			break

		}
	}

	if flag == false {
		fmt.Println("ERROR!")
	} else {
		fmt.Println("Blocks verified. No Changing...")
	}
	return flag
}

// (7) Function for calculating hash of a transaction or a block. If the size of the transaction is very
// large, then Merkle‐Damgard Transform shall be used.
func CalculateHash(chain *blockchain) {

	for i := 0; i < len(chain.list); i++ {
		T_Trx := TraverseTillindex(&chain.list[i].Arr, &chain.list[i].M_index)
		var attr string
		attr += strconv.Itoa(chain.list[i].id)
		attr += T_Trx + chain.list[i].Prev_hash
		sum := sha256.Sum256([]byte(attr))
		chain.list[i].Curr_hash = fmt.Sprintf("%x", sum)
		if i < len(chain.list)-1 {
			chain.list[i+1].Prev_hash = fmt.Sprintf("%x", sum)
		}

	}
}

func (blocklist *blockchain) AddBlock(x int) *block {
	tempblock := NewBlock(x)

	if VerifyChain(blocklist) {
		blocklist.list = append(blocklist.list, tempblock)
		CalculateHash(blocklist)

		fmt.Println("Block added.")
	} else {
		fmt.Println(" ERROR! Block can not be added")
		return nil
	}
	return tempblock
}

func main() {

	chain := new(blockchain)
	var x = 0
	for i := 0; i < 10; i++ {
		chain.AddBlock(x + i)
	}

	println("Mining... ")
	Mineblock(chain)

	//Display Blockchain
	DisplayBlocks(chain)

	fmt.Println("Changing block with id:", 0)
	ChangeBlock(chain, 0)

	//Display Blockchain
	DisplayBlocks(chain)

}

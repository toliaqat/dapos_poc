package main

import (
	"math/rand"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func logSeparator() {
	log.Info("~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~ ~~~~")
}
func prefixLinesWith(lines []string, prefix string) []string {
	var prefixedLines = []string{}

	for index, line := range lines {
		if index == 0 {
			prefixedLines = append(prefixedLines, line)
		} else {
			prefixedLines = append(prefixedLines, prefix+line)
		}
	}

	return prefixedLines
}

func GetRandomNumber(boundary int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(boundary)
}

func getDictKeysAsList() []string {
	keys := make([]string, 0)
	for k, _ := range getNodes() {
		keys = append(keys, k)
	}

	return keys
}

func getRandomNode(nodeToIgnore *Node) *Node {
	nodes := getNodes()
	nodesNames := getDictKeysAsList()

	var theNode *Node
	randomNum := GetRandomNumber(len(nodes))
	theNode = getNodes()[nodesNames[randomNum]]

	return theNode
}

func sendRandomTransaction(fromWallet string, toWallet string, transactionId int, amount int, delegate *Node) {

	transaction := Transaction{
		transactionId,
		fromWallet,
		toWallet,
		amount,
		time.Now(),
	}

	go func() {
		log.Infof("sendRandomTx()  | Tx_%d(%s -> %s) | send to delegate -> %s",
			transaction.Id,
			transaction.From,
			transaction.To,
			delegate.Wallet,
		)
		delegate.TxChannel <- transaction
	}()
}

var nodes *map[string]*Node

var once sync.Once

var startingBalance = 10000000

func CreateNodeAndAddToList(newMember string) {

	GenesisBlock := Block{
		Prev: nil,
		Next: nil,
		Transaction: &Transaction{
			0,
			"dl",
			"dl",
			startingBalance,
			time.Now(),
		},
	}

	node := Node{
		GenesisBlock:    	GenesisBlock,
		LastBlock:       	&GenesisBlock,
		TxChannel:       	make(chan Transaction),
		TxCount:			0,
		Wallet:       		string(newMember),
		IsDelegate: 		true,
		TxFromChainById:	map[int]*Transaction{},
		AllWallets:			map[string]int{},
	}

	node.AllWallets["dl"] = startingBalance
	getNodes()[newMember] = &node
	node.StartListenForTx()
}

func getNodes() map[string]*Node {
	once.Do(func() {
		nodes = &map[string]*Node{}
	})
	return *nodes
}

func getNodeByAddress(address string) *Node {
	var theNode = getNodes()[string(address)]
	return theNode
}

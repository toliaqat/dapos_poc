package main

import (
	"fmt"

	"time"
)

func (node *Node) StartListenForTx() {
	//gets initiated on node initialization
	node.LogLines = append(node.LogLines, fmt.Sprintf("StartListenForTx - Delegate - %s", node.Wallet))

	go func() {
		for {
			tx := <-node.TxChannel

			var additionalLogLines = []string{}

			node.LogLines = append(node.LogLines, fmt.Sprintf("    GotTX()-node    | Tx_%d(From:%s -> To:%s) | CurrentNode:%s | RecievedFrom:%s", tx.Id, tx.From, tx.To, node.Wallet,tx.DelId))

			seen := node.checkIfValidated(tx.Id)

			if seen { //if the tx has already been validated, log and do nothing
				node.LogLines = append(node.LogLines, fmt.Sprintf("        delegate %s: skipping received transaction %d", node.Wallet, tx.Id))
			} else { //else check tx for validity
				additionalLogLines = node.validateBlockAndTransmit(&tx)
			}

			additionalLogLines = prefixLinesWith(additionalLogLines, "        ", "            ")

			node.LogLines = append(node.LogLines, additionalLogLines...)

			TotalTxProcessed++
		}
	}()
}

func (node *Node) checkIfValidated(txId int) bool {
	if node.TxFromChainById[txId] == nil {
		return false
	} else {
		return true
	}
}

func (node *Node) validateBlockAndTransmit(tx *Transaction) []string {
	var logLines = []string{}
	var additionalLogLines = []string{}

	logLines = append(logLines, fmt.Sprintf("validateBlock()"))

	//call Validate(transaction)
	valid := node.validate(tx)

	additionalLogLines = prefixLinesWith(additionalLogLines, "", "    ")
	logLines = append(logLines, additionalLogLines...)


	if valid {
		logLines = append(logLines, fmt.Sprintf("Node ID: %s, Transaction ID: %d, Value: %d", node.Wallet, tx.Id, tx.Value))
		logLines = append(logLines, fmt.Sprintf("delegate %s: received valid transaction %d with value: %d", node.Wallet, tx.Id, tx.Value))

		//add valid transaction to 'validated' list
		node.TxFromChainById[tx.Id] = tx


		//report back if no more expected tx
		//if tx was last expected (4) then report balances
		node.TxCount++
		if node.TxCount == 1 {
			node.StartTime = time.Now()
		}

		if node.TxCount >= (NrOfTx) {
			logLines = append(logLines, fmt.Sprintf("Node %s thinks balance of BobSt: %d, Chris: %d, GregM: %d, Muham: %d \n",
				node.Wallet,
				node.AllWallets["BobSt"],
				node.AllWallets["Chris"],
				node.AllWallets["GregM"],
				node.AllWallets["Muham"]))

			TimeToComplete := time.Since(node.StartTime)

			logLines = append(logLines, fmt.Sprintf("Delegate %s processed %d transactions in %d time", node.Wallet, 4, TimeToComplete))

		}

		// set the delegate id to current id and broadcast the valid transaction to other nodes
		for k, _ := range getNodes() {
			destinationNode := getNodes()[k]
			if destinationNode.Wallet == node.Wallet || destinationNode.Wallet == tx.DelId{
				continue
			}

			go func() {
				tx.DelId = node.Wallet
				destinationNode.TxChannel <- *tx
			}()
			//TODO: SendingFrom not correctly printing (tx.DelId should be = to node.Wallet, instead printing old value)
			logLines = append(logLines, fmt.Sprintf("sendTx()        | Tx_%d(From:%s -> To:%s) | CurrentNode:%s |SendingFrom%s -> SendingTo:%s", tx.Id, tx.From, tx.To, node.Wallet,tx.DelId, destinationNode.Wallet))
		}
	} else {
		logLines = append(logLines, fmt.Sprintf("Node ID: %s, Transaction: %d, From ID: %s, Value: %d", node.Wallet, tx.Id, tx.From, tx.Value))
		logLines = append(logLines, fmt.Sprintf("delegate %s: received invalid transaction %d with value: %d", node.Wallet, tx.Id, tx.Value))
	}

	return logLines
}

//validates the transaction and adds it to the end of the chain
func (node *Node) validate(tx *Transaction) bool {
	//don't process a negative tx
	if tx.Value < 0 {
		return false
	}

	//check if transaction goes at end of list, then AllWallets can check validity
	//if tx.Time.After(node.LastBlock.Transaction.Time) {
	if true {

		if node.AllWallets[tx.From] < tx.Value { //sender doesn't have enough money
			return false
		} else { //transaction is valid!!!
			//update AllWallets balance
			node.AllWallets[tx.From] -= tx.Value
			node.AllWallets[tx.To] += tx.Value

			//add tx to end of list
			node.LastBlock.Next = &Block{
				node.LastBlock,
				nil,
				tx,
			}
			node.LastBlock = node.LastBlock.Next

			//return true then add to TxFromChainById & broadcast to delegates
			return true
		}
	} else { //if tx is not at end of list, iterate backwards to find balances of time of tx
		//TODO: support tx that come before lastBlock

		//start with node.AllWallets and iterate backwards from node.LastBlock until
		//the state of the chain at time of transaction is discovered
		/*
		Node GregM thinks balance of BobSt: 998, Chris: 1003, GregM: 1000, Muham: 999
		Node Muham thinks balance of BobSt: 998, Chris: 1003, GregM: 1000, Muham: 999
		Node BobSt thinks balance of BobSt: 998, Chris: 1003, GregM: 1000, Muham: 999
		Node Chris thinks balance of BobSt: 998, Chris: 1003, GregM: 1000, Muham: 999
		*/

	}

	return false
}

func (node *Node) DumpLogLines() {
	for _, line := range node.LogLines {
		fmt.Println(line)
	}
}

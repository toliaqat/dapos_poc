package main

import (
	"log"
)

type VoteCounter struct {
	votes	  	map[int]*Votes
	Channel 	chan Vote
	quit        chan bool
}

type Votes struct {
	Transaction  Transaction
	VoteYesNo	 []bool
	VoteCount    int
	NbrDelegates int
}

type Vote struct {
	TransactionId int
	VoteYesNo	  bool
	DelegateId    int
}

func NewVoteCounter(c chan Vote) *VoteCounter {
	vm := make(map[int]*Votes)
	return &VoteCounter{
		votes: vm,
		Channel: c,
	}
}

func (vc *VoteCounter)AddVoting(t Transaction, nbrDelegates int) {
	votes := Votes{
		Transaction: 	t,
		VoteYesNo:      make([]bool, 0),
		VoteCount:   	0,
		NbrDelegates: 	nbrDelegates,
	}
	vc.votes[t.Id] = &votes
}

func (vc *VoteCounter) Start() {
	go func() {
		for {
			select {
			case vote := <-vc.Channel:
				// we have received a vote.
				v := vc.votes[vote.TransactionId]
				log.Printf("Received Vote for transaction: %d value %t from delegate %d with value %d", vote.TransactionId, vote.VoteYesNo, vote.DelegateId, v.Transaction.Value)
				v.VoteYesNo = append(v.VoteYesNo, vote.VoteYesNo)
				v.VoteCount++
				log.Printf("nbr delegates = %d and nbr votes = %d", v.NbrDelegates, v.VoteCount)
				if(v.NbrDelegates == v.VoteCount) {
					if v.isValid() {
						updateAccounts(v.Transaction)
					} else {
						log.Printf("The Delegates voted this transaction as invalid %v", v.Transaction )
					}
				}
			case <-vc.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

func (v Votes) isValid() bool {
	var positiveCount = 0
	var negativeCount = 0
	for _, value := range v.VoteYesNo {
		if(value) {
			positiveCount++
		} else {
			negativeCount++
		}
	}
	if (positiveCount == (v.NbrDelegates)) {
		return true
	} else {
		return false
	}
}

func updateAccounts(t Transaction) {
	log.Printf("Update Accounts: %d", t.Id)
	fromAcct := GetAccount(t.From)
	toAcct := GetAccount(t.To)
	fromAcct.Balance -= t.Value
	toAcct.Balance += t.Value
	fromAcct.Transactions = append(fromAcct.Transactions, t)
	toAcct.Transactions = append(toAcct.Transactions, t)
}

// Stop signals the worker to stop listening for work requests.
func (vc VoteCounter) Stop() {
	go func() {
		vc.quit <- true
	}()
}
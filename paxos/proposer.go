package paxos

import (
	"log"
	"time"
)

func NewProposer(client Client) *Proposer {
	return &Proposer{client: client}
}

type Proposer struct {
	client    Client
	lasMsgSeq int
	acceptors []Acceptor
	proposal  Proposal
}

func (p *Proposer) Run(value string) {
	p.newProposal(value)
	p.broadcast(Prepare)

	promiseCnt := 0
	for {
		msg, err := p.client.Recv(time.Second)
		if err != nil {
			log.Printf("proposer: %d failed to recive msg. err is %s", p.client.GetId(), err.Error())
			continue
		}

		switch msg.GetType() {
		case Promise:
			promiseCnt++
			if msg.GetProposal() != (Proposal{}) {
				p.proposal.value = msg.GetProposal().value
			}
			if p.reachMajority(promiseCnt) {
				p.broadcast(Accept)
			}
		default:
			log.Panicf("acceptor: %d unexpected message type: %v", a.client.GetId(), msg.typ)
		}
	}
}

func (p *Proposer) reachMajority(i int) bool {
	return i > len(p.acceptors)/2+1
}

func (p *Proposer) newProposal(value string)  {
	p.lasMsgSeq++
	p.proposal = Proposal{id: p.lasMsgSeq<<16 | p.client.GetId(), value: value}
}

func (p *Proposer) broadcast(typ MessageType) {
	for _, acceptor := range p.acceptors {
		err := p.client.Send(NewMessage(p.client.GetId(), acceptor.GetClientId(), p.proposal.GetId(), typ, p.proposal))
		if err != nil {
			log.Fatal(err)
		}
	}
}

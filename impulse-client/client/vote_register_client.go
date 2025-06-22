package client

import (
	"github.com/andantan/vote-blockchain-server/impulse-client/data"
)

func RegisterUsers() {
	cfg := data.GetRegisterData()

	p := NewProposalClient(1)
	p.RequestProposal(cfg.Votes[0])
	p.Wg.Wait()

	s := NewSubmitClient(1, nil)
	s.RequestSubmitLoop(cfg.Votes[0])
	s.Wg.Wait()
}

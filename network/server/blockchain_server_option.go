package server

import (
	"time"

	"github.com/andantan/vote-blockchain-server/network/server/listener"
)

type ListenerOption struct {
	*listener.VoteProposalListenerOption
	*listener.VoteSubmitListenerOption
}

func NewListenerOption() ListenerOption {
	return ListenerOption{}
}

func (o *ListenerOption) SetVoteProposalListenerOption(network string, port, channelBufferSize uint16) {
	o.VoteProposalListenerOption = listener.NewVoteProposalListenerOption(network, port, channelBufferSize)
}

func (o *ListenerOption) SetVoteSubmitListenerOption(network string, port, channelBufferSize uint16) {
	o.VoteSubmitListenerOption = listener.NewVoteSubmitListenerOption(network, port, channelBufferSize)
}

type BlockOption struct {
	blockTime time.Duration
	maxTxSize uint32
}

func NewBlockOption(blockTime time.Duration, maxTxSize uint32) BlockOption {
	return BlockOption{
		blockTime: blockTime,
		maxTxSize: maxTxSize,
	}
}

type StoreOption struct {
	StoreBaseDirectory   string
	StoreBlocksDirectory string
}

func NewStoreOption(baseDir, blocksDir string) StoreOption {
	return StoreOption{
		StoreBaseDirectory:   baseDir,
		StoreBlocksDirectory: blocksDir,
	}
}

type BlockChainServerOpts struct {
	ListenerOption
	BlockOption
	StoreOption
}

func NewBlockChainServerOpts() BlockChainServerOpts {
	return BlockChainServerOpts{}
}

func (o *BlockChainServerOpts) SetListenerOption(option ListenerOption) {
	o.ListenerOption = option
}

func (o *BlockChainServerOpts) SetBlockOption(option BlockOption) {
	o.BlockOption = option
}
func (o *BlockChainServerOpts) SetStoreOption(option StoreOption) {
	o.StoreOption = option
}

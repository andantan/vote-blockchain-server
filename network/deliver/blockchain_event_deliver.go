package deliver

import (
	"github.com/andantan/vote-blockchain-server/core/block"
	"github.com/andantan/vote-blockchain-server/core/mempool"
	"github.com/andantan/vote-blockchain-server/network/deliver/event"
)

type EventDeliver struct {
	defaultProtocol            string
	defaultAddress             string
	defaultPort                uint16
	CreatedBlockEventDeliver   *event.CreatedBlockeventUnicaster
	ExpiredPendingEventDeliver *event.ExpiredPendingEventUnicaster
}

func NewEventDeliver(defaultProtocol, defaultAddress string, defaultPort uint16) *EventDeliver {
	return &EventDeliver{
		defaultProtocol: defaultProtocol,
		defaultAddress:  defaultAddress,
		defaultPort:     defaultPort,
	}
}

func (d *EventDeliver) SetCreatedBlockEventDeliver(path string) {
	cfg := event.NewCreatedBlockEventEndpoint(d.defaultProtocol, d.defaultAddress, d.defaultPort, path)
	d.CreatedBlockEventDeliver = event.NewCreatedBlockeventUnicaster(cfg)
}

func (d *EventDeliver) UnicastCreatedBlockEvent(blk *block.Block) {
	d.CreatedBlockEventDeliver.Unicast(blk)
}

func (d *EventDeliver) SetExpirdPendingEventDeliver(path string) {
	cfg := event.NewExpiredPendingEventEndPoint(d.defaultProtocol, d.defaultAddress, d.defaultPort, path)
	d.ExpiredPendingEventDeliver = event.NewExpiredPendingEventUnicaster(cfg)
}

func (d *EventDeliver) UnicastExpiredPendingEvent(expiredPended *mempool.Pended) {
	d.ExpiredPendingEventDeliver.Unicast(expiredPended)
}

package socket

import (
	"sync"
	"time"

	"github.com/zishang520/engine.io/events"
	"github.com/zishang520/engine.io/packet"
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/socket.io/parser"
)

type SocketId string

type Room string

type WriteOptions struct {
	packet.Options

	Volatile     bool
	PreEncoded   bool
	WsPreEncoded string
}

type BroadcastFlags struct {
	WriteOptions

	Local     bool
	Broadcast bool
	Binary    bool
	Timeout   *time.Duration
}

type BroadcastOptions struct {
	Rooms  *types.Set[Room]
	Except *types.Set[Room]
	Flags  *BroadcastFlags
}

type Adapter interface {
	New(NamespaceInterface) Adapter

	Rooms() *sync.Map
	Sids() *sync.Map
	Nsp() NamespaceInterface

	// To be overridden
	Init()

	// To be overridden
	Close()

	// Returns the number of Socket.IO servers in the cluster
	ServerCount() int64

	// Adds a socket to a list of room.
	AddAll(SocketId, *types.Set[Room])

	// Removes a socket from a room.
	Del(SocketId, Room)

	// Removes a socket from all rooms it's joined.
	DelAll(SocketId)

	SetBroadcast(func(*parser.Packet, *BroadcastOptions))
	// Broadcasts a packet.
	//
	// Options:
	//  - `Flags` {*BroadcastFlags} flags for this packet
	//  - `Except` {*types.Set[Room]} sids that should be excluded
	//  - `Rooms` {*types.Set[Room]} list of rooms to broadcast to
	Broadcast(*parser.Packet, *BroadcastOptions)

	// Broadcasts a packet and expects multiple acknowledgements.
	//
	// Options:
	//  - `Flags` {*BroadcastFlags} flags for this packet
	//  - `Except` {*types.Set[Room]} sids that should be excluded
	//  - `Rooms` {*types.Set[Room]} list of rooms to broadcast to
	BroadcastWithAck(*parser.Packet, *BroadcastOptions, func(uint64), func(...any))

	// Gets a list of sockets by sid.
	Sockets(*types.Set[Room]) *types.Set[SocketId]

	// Gets the list of rooms a given socket has joined.
	SocketRooms(SocketId) *types.Set[Room]

	// Returns the matching socket instances
	FetchSockets(*BroadcastOptions) []any

	// Makes the matching socket instances join the specified rooms
	AddSockets(*BroadcastOptions, []Room)

	// Makes the matching socket instances leave the specified rooms
	DelSockets(*BroadcastOptions, []Room)

	// Makes the matching socket instances disconnect
	DisconnectSockets(*BroadcastOptions, bool)

	// Send a packet to the other Socket.IO servers in the cluster
	ServerSideEmit(string, ...any) error
}

type SocketDetails interface {
	Id() SocketId
	Handshake() *Handshake
	Rooms() *types.Set[Room]
	Data() any
}

type NamespaceInterface interface {
	EventEmitter() *StrictEventEmitter

	On(string, ...events.Listener) error
	Once(string, ...events.Listener) error
	EmitReserved(string, ...any)
	EmitUntyped(string, ...any)
	Listeners(string) []events.Listener

	Sockets() *sync.Map
	Server() *Server
	Adapter() Adapter
	Name() string
	Ids() uint64

	// Sets up namespace middleware.
	Use(func(*Socket, func(*ExtendedError))) NamespaceInterface

	// Targets a room when emitting.
	To(...Room) *BroadcastOperator

	// Targets a room when emitting.
	In(...Room) *BroadcastOperator

	// Excludes a room when emitting.
	Except(...Room) *BroadcastOperator

	// Adds a new client.
	Add(*Client, any, func(*Socket)) *Socket

	// Emits to all clients.
	Emit(string, ...any) error

	// Sends a `message` event to all clients.
	Send(...any) NamespaceInterface

	// Sends a `message` event to all clients.
	Write(...any) NamespaceInterface

	// Emit a packet to other Socket.IO servers
	ServerSideEmit(string, ...any) error

	// Gets a list of clients.
	AllSockets() (*types.Set[SocketId], error)

	// Sets the compress flag.
	Compress(bool) *BroadcastOperator

	// Sets a modifier for a subsequent event emission that the event data may be lost if the client is not ready to
	// receive messages (because of network slowness or other issues, or because they’re connected through long polling
	// and is in the middle of a request-response cycle).
	Volatile() *BroadcastOperator

	// Sets a modifier for a subsequent event emission that the event data will only be broadcast to the current node.
	Local() *BroadcastOperator

	// Adds a timeout in milliseconds for the next operation
	//
	// <pre><code>
	//
	// io.Timeout(1000 * time.Millisecond).Emit("some-event", func(args ...any) {
	//   // ...
	// });
	//
	// </pre></code>
	Timeout(time.Duration) *BroadcastOperator

	// Returns the matching socket instances
	FetchSockets() ([]*RemoteSocket, error)

	// Makes the matching socket instances join the specified rooms
	SocketsJoin(...Room)

	// Makes the matching socket instances leave the specified rooms
	SocketsLeave(...Room)

	// Makes the matching socket instances disconnect
	DisconnectSockets(bool)
}

// Package trans contains tools allowing inter-agent resource exchange.
package trans

import (
	"errors"
	"github.com/rwcarlsen/goclus/rsrc"
)

// TransType indicates a transaction's type (e.g. offer or request).
type TransType int

const (
	// Offer indicates a removal of resources from the transaction supplier.
	Offer TransType = iota
	// Request indicates addition of resources to the transaction requester.
	Request
)

var listeners []Listener

// Listener is implemented by entities that desire to receive notifications
// every time a transaction is approved and executed between any two
// simulation agents.
type Listener interface {
	TransNotify(*Transaction)
}

// ListenAll adds l to a global list of agents that receive notifications
// for every approved transaction (usually used by "special" agents e.g.
// book-keeper, etc.).
// These notifications are sent when the Approve method is called - 
// before Approve returns and directly after the resource transfer.
// Simulation execution continues only after l's TransNotify method returns.
func ListenAll(l Listener) {
	listeners = append(listeners, l)
}

func notifyListeners(t *Transaction) {
	for _, l := range listeners {
		l.TransNotify(t)
	}
}

// Supplier is implemented by all agents that are able to send resources to
// other agents via matched/approved transactions.
type Supplier interface {
	RemoveResource(*Transaction)
}

// Requester is implemented by all agents that are able to receive
// resources from other agents via matched/approved transactions.
type Requester interface {
	AddResource(*Transaction)
}

// Transaction allows agents to inform each other about desired resource
// exchange.
// Partial transactions (offers/requests) are generally sent as part of a
// message and matched by transaction-matching agents. The matched
// transaction is returned to the supplier who then (generally) calls the
// Approve method to initiate the resource transfer.
type Transaction struct {
	tp       TransType
	res      rsrc.Resource
	Sup      Supplier
	Req      Requester
	Manifest []rsrc.Resource
}

// NewOffer creates a new offer transaction.
func NewOffer(sup Supplier) *Transaction {
	return &Transaction{
		tp:  Offer,
		Sup: sup,
	}
}

// NewOffer creates a new request transaction.
func NewRequest(req Requester) *Transaction {
	return &Transaction{
		tp:  Request,
		Req: req,
	}
}

// MatchWith pairs a set of offer-request transactions by setting the
// offer's requester to the request's requester and the request's supplier
// to the supplier's supplier.
func (t *Transaction) MatchWith(other *Transaction) error {
	if t.tp == other.tp {
		return errors.New("trans: Non-complementary transaction types")
	}

	if t.tp == Offer {
		t.Req = other.Req
		other.Sup = t.Sup
	} else {
		t.Sup = other.Sup
		other.Req = t.Req
	}
	return nil
}

// Type returns the type (Offer or Request) of the transaction.
func (t *Transaction) Type() TransType {
	return t.tp
}

// Approve executes the resource transfer: resources are removed from the
// supplier and given to the requester.
// All transaction notification listeners are also notified immediately
// following the resource transfer before Approve returns.
func (t *Transaction) Approve() {
	t.Sup.RemoveResource(t)
	t.Req.AddResource(t)
	notifyListeners(t)
}

// Resource returns the resource associated with this transaction (not a
// clone).
func (t *Transaction) Resource() rsrc.Resource {
	return t.res
}

// SetResource sets this transaction's resource to a clone of r.
func (t *Transaction) SetResource(r rsrc.Resource) {
	t.res = r.Clone()
}

// Clone returns a shallow copy of the transaction except the
// the resource of the returned transaction is a clone of the original
// resource.
func (t *Transaction) Clone() *Transaction {
	clone := *t
	clone.res = t.res.Clone()
	return &clone
}

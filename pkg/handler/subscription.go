package handler

import (
	"encoding/json"
	"fmt"
	"github.com/n3wscott/cloudevents-discovery/pkg/apis/subscription"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	once          sync.Once
	subscriptions map[string]subscription.Subscription
}

// TODO: I made a choice to not implement the OpenAPI of the current api for subscription. I wanted id in the url, not query.

func (h *SubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(func() {
		h.subscriptions = make(map[string]subscription.Subscription, 0)

		subscriptions := make([]subscription.Subscription, 0)
		err := json.Unmarshal([]byte(exampleSubscriptions), &subscriptions)
		if err != nil {
			panic(err)
		}

		for _, sub := range subscriptions {
			h.subscriptions[sub.ID] = sub
		}
	})

	vars := mux.Vars(r)
	id := vars["id"]

	switch r.Method {
	case http.MethodOptions:
		if id != "" {
			w.Header().Set("Allow", "GET,OPTIONS")
		} else {
			w.Header().Set("Allow", "GET,PUT,POST,DELETE,OPTIONS")
		}

	case http.MethodGet:
		if id == "" {
			h.handleQuery(w, r)
		} else {
			h.handleGet(id, w, r)
		}

	case http.MethodPost, http.MethodPut:
		if id != "" {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}
		h.handleCreateOrUpdate(w, r)

	case http.MethodDelete:
		h.handleDelete(id, w, r)

	default:
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// 3.2.4.1. Creating a subscription
// The Create operation SHOULD be supported by compliant Event Producers. It creates a new Subscription. The client proposes a subscription object which MUST contain all REQUIRED properties. The subscription manager then realizes the subscription and returns a subscription object that also contains all OPTIONAL properties for which default values have been applied.
//
// Parameters:
//
// subscription (subscription) - REQUIRED. Proposed subscription object.
// Result:
//
// subscription (subscription) - REQUIRED. Realized subscription object.
// Errors:
//
// ok - the operation succeeded
// conflict - a subscription with the given id already exists
// invalid - the proposed subscription object contains invalid information
// Protocol bindings MAY map the Create operation such that the proposed id is ignored and the subscription manager assigns one instead.
//
// 3.2.4.4. Updating a Subscription
// The Update operation MAY be supported by compliant Event Producers. To request the update of a Subscription, the client submits a proposed subscription object whose id MUST match an existing subscription. All other properties MAY differ from the original subscription. The subscription manager then updates the subscription and returns a subscription object that also contains all OPTIONAL properties for which default values have been applied.
//
// Parameters:
//
// subscription (subscription) - REQUIRED. Proposed subscription object.
// Result:
//
// subscription (subscription) - REQUIRED. Realized subscription object.
// Protocol bindings MAY map the Update and the Create operation into a composite "upsert" operation that creates a new subscription if one with the given id does not exist. In this case, the operation is *Create and follows that operation's rules.
func (h *SubscriptionHandler) handleCreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	sub := new(subscription.Subscription)

	err := json.NewDecoder(r.Body).Decode(sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, found := h.subscriptions[sub.ID]; found && r.Method == http.MethodPost {
		http.Error(w, fmt.Sprintf("subscription %q already exists", sub.ID), http.StatusConflict)
		return
	}

	// TODO: validate all of the subscription.
	if sub.ID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Save.
	h.subscriptions[sub.ID] = *sub

	js, err := json.Marshal(sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// 3.2.4.2. Retrieving a Subscription
// The Retrieve operation MUST be supported by compliant Event Producers. It returns the specification of the identified subscription.
//
// Parameters:
//
// id (string) - REQUIRED. Identifier of the subscription.
// Result:
//
// subscription (subscription) - REQUIRED. Subscription object.
// Errors:
//
// ok - the operation succeeded
// notfound - a subscription with the given id already exists
func (h *SubscriptionHandler) handleGet(id string, w http.ResponseWriter, r *http.Request) {
	var sub *subscription.Subscription

	if s, ok := h.subscriptions[id]; ok {
		sub = &s
	}

	if sub == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	js, err := json.Marshal(sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// 3.2.4.3. Querying for a list of Subscriptions
// The Query operation SHOULD be supported by compliant Event Producers. It allows to query the list of subscriptions on the subscription manager associated with or otherwise visible to the party making the request. If supported, it MUST be supported at the same endpoint as the Create subscription operation.
//
// Parameters:
//
// none
// Result:
//
// subscription (list of subscription) - REQUIRED. List of subscription objects
// Errors:
//
// ok - the operation succeeded and returned results
// nocontent - the operation succeeded and returned no results
// Protocol bindings and implementations of such bindings MAY add custom filter constraints and pagination arguments as parameters. A request without filtering constraints SHOULD return all available subscriptions associated with or otherwise visible to the party making the request.
func (h *SubscriptionHandler) handleQuery(w http.ResponseWriter, r *http.Request) {
	subscriptions := make([]subscription.Subscription, 0)

	// collect
	for _, sub := range h.subscriptions {
		subscriptions = append(subscriptions, sub)
	}

	js, err := json.Marshal(subscriptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

// 3.2.8. Deleting a Subscription
// The Delete operation SHOULD be supported by compliant Event Producers. It returns the specification of the identified subscription.
//
// Parameters:
//
// id (string) - REQUIRED. Identifier of the subscription.
// Result:
//
// subscription (subscription) - REQUIRED. Subscription object.
// Errors:
//
// ok - the operation succeeded
// notfound - a subscription with the given id already exists // TODO: fix this in the spec upstream. Should be: _does not_ already exist.
func (h *SubscriptionHandler) handleDelete(id string, w http.ResponseWriter, r *http.Request) {
	if id == "" {
		http.Error(w, fmt.Sprintf("subscription %q not found", id), http.StatusNotFound)
		return
	}

	if _, found := h.subscriptions[id]; !found {
		http.Error(w, fmt.Sprintf("subscription %q not found", id), http.StatusNotFound)
		return
	}

	delete(h.subscriptions, id)
	w.WriteHeader(http.StatusOK)
}

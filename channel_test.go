// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE file.

package channel_subscription

import (
	"fmt"
	"os"
	"testing"
)

var cs *ChannelSubscription
var numChs = 10
var channels []chan struct{}
var ret chan struct{}

func TestMain(m *testing.M) {
	cs = NewChannelSubscription()
	ret = make(chan struct{}, numChs)
	channels = make([]chan struct{}, numChs)

	for i := 0; i < numChs; i++ {
		channels[i] = make(chan struct{})
		go func(ch chan struct{}) {
			for {
				ret <- <-ch
			}
		}(channels[i])
	}

	os.Exit(m.Run())
}

func TestSubscribe(t *testing.T) {
	for _, ch := range channels {
		cs.Subscribe(ch)
	}
}

func TestSubscribeFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Subscribe don't failed.")
		}
	}()
	cs.Subscribe(0)
}

func TestUnsubscribe(t *testing.T) {
	ch := make(chan bool)

	err := cs.Unsubscribe(ch)
	if err != nil && err != ErrNotFound {
		t.Fatal(err)
	}

	cs.Subscribe(ch)

	err = cs.Unsubscribe(ch)
	if err != nil {
		t.Fatal(err)
	}
}

func countRecv(t *testing.T) {
	count := 0
	for i := 0; i < numChs; i++ {
		if _, ok := <-ret; ok {
			count++
		}
	}
	if count != numChs {
		t.Fatal("wrong number of answers", count)
	}
}

func TestTrySend(t *testing.T) {
	ok := cs.TrySend(struct{}{})
	if !ok {
		t.Fatal("TrySend failed")
	}
	countRecv(t)
}

func TestSend(t *testing.T) {
	cs.Send(struct{}{})
	countRecv(t)
}

func Example() {
	const numChs = 10
	cs := NewChannelSubscription()
	ret := make(chan bool, numChs)
	channels := make([]chan bool, numChs)

	for i := 0; i < numChs; i++ {
		channels[i] = make(chan bool)
		go func(ch chan bool) {
			for {
				ret <- <-ch
			}
		}(channels[i])
		cs.Subscribe(channels[i])
	}

	cs.Send(true)

	count := 0
	for i := 0; i < numChs; i++ {
		if <-ret {
			count++
		}
	}
	if count != numChs {
		fmt.Println("Send failed.")
	} else {
		fmt.Println("Send ok!")
	}
	// Output: Send ok!
}

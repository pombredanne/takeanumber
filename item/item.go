// Copyright 2015 Daniel Lindsley. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package item implements a unit of data (+ metadata) for use in the Queue.
package item

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"strings"
	"time"
)

// An error for when the provided body is empty.
var EmptyBody = errors.New("No body provided.")

// The Item itself.
type Item struct {
	Id               string
	Body             string
	InitialRetries   int
	RemainingRetries int
	Reserved         bool
	Created          time.Time
}

// Decrements the number of times the Item can be retried.
//
// If successfully decremented, this will return true. If there are no retries
// left, this will return false.
func (i *Item) DecrRetries() bool {
	if !i.ShouldRetry() {
		return false
	}

	i.RemainingRetries--
	return true
}

// Returns if the Item can be retried.
//
// Returns true if it can be retried, false if all retries have been used.
func (i *Item) ShouldRetry() bool {
	return i.RemainingRetries > 0
}

// Returns if the Item is reserved.
//
// Returns true if reserved, false if not.
func (i *Item) IsReserved() bool {
	return i.Reserved
}

// Marks an Item as reserved.
func (i *Item) Reserve() {
	i.Reserved = true
}

// Releases the reserved status on an Item.
func (i *Item) Release() {
	i.Reserved = false
}

// New creates a new Item instance.
//
// If an empty body is provided, this will return an EmptyBody error.
func New(body string, retries int) (*Item, error) {
	if len(strings.TrimSpace(body)) <= 0 {
		return &Item{}, EmptyBody
	}

	id := uuid.New()
	created := time.Now()
	return &Item{id, body, retries, retries, false, created}, nil
}

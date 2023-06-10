// This file contains types inside the asapp domain.

// Package app contains structs and implementations for this app.
package app

import (
	"github.com/Bitspark/go-bitnode/bitnode"
	"log"
)

// @@TYPES@@

// DOMAIN STRUCT

// Domain containing mainly wrappers for applications.
type Domain struct {
	Domain *bitnode.Domain
	Node   bitnode.Node
}

// NewBlankSparkable creates a new BlankSparkable instance.
func (asapp *Domain) NewBlankSparkable() (*BlankSparkable, error) {
	// Get the BlankSparkable sparkable from the domain.
	blankSblSpark, err := asapp.Domain.GetSparkable("fullBlankDomain.BlankSparkable")
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the BlankSparkable spark.
	blankSblSpk, err := asapp.Node.PrepareSystem(bitnode.Credentials{}, *blankSblSpark)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the BlankSparkable.
	blankSbl := &BlankSparkable{
		System: blankSblSpk,
	}

	return blankSbl, nil
}

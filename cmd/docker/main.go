package main

import (
	"blank-sparkable/app"
	"context"
	"github.com/Bitspark/go-bitnode/api/wsApi"
	"github.com/Bitspark/go-bitnode/bitnode"
	"github.com/Bitspark/go-bitnode/factories"
	"github.com/Bitspark/go-bitnode/store"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	localAddress := os.Getenv("BITNODE_LOCAL_ADDRESS")
	remoteNodeAddress := os.Getenv("BITNODE_REMOTE_ADDRESS")

	node := bitnode.NewNode()
	node.AddMiddlewares(factories.GetMiddlewares())

	dom := bitnode.NewDomain()
	dom, _ = dom.AddDomain("hub")

	// Prepare node connections.
	nodeConns := wsApi.NewNodeConns(node, remoteNodeAddress)

	// Add factories.
	if err := node.AddFactory(wsApi.NewWSFactory(nodeConns)); err != nil {
		log.Fatal(err)
	}

	// Prepare node.
	if err := dom.LoadFromDir("./domain", true); err != nil {
		log.Fatal(err)
	}
	if err := dom.Compile(); err != nil {
		log.Fatal(err)
	}

	blankDomain := &app.Domain{
		Domain: dom,
		Node:   node,
	}

	// Read store.
	st1 := store.NewStore("store")
	if err := st1.Read("."); err != nil {
		log.Println(err)
	} else {
		// Load node.
		if err := nodeConns.Load(st1, dom); err != nil {
			log.Fatalf("Error loading node: %v", err)
		} else {
			log.Printf("Loaded node from %s", ".")
		}
	}

	creds := bitnode.Credentials{}

	var blankSbl *app.BlankSparkable

	created := false
	if len(node.Systems(creds)) == 0 {
		var err error
		blankSbl, err = blankDomain.NewBlankSparkable()
		if err != nil {
			log.Fatal(err)
		}

		// Make computer system the root system.
		node.SetSystem(blankSbl.Native())

		created = true
	} else {
		log.Printf("Found %d startup systems", len(node.Systems(creds)))

		// Get the system from the node.
		blankSblSys := node.System(creds)

		blankSbl = &app.BlankSparkable{
			System: blankSblSys,
		}
	}

	// Add the custom BlankSparkable implementation.
	if err := blankSbl.Init(); err != nil {
		log.Fatal(err)
	}

	// Handle loading.
	if created {
		blankSbl.AddCallback(bitnode.LifecycleCreate, bitnode.NewNativeEvent(func(vals ...bitnode.HubItem) error {
			return blankSbl.Native().EmitEvent(bitnode.LifecycleLoad)
		}))
	} else {
		if err := blankSbl.Native().EmitEvent(bitnode.LifecycleLoad); err != nil {
			log.Fatal(err)
		}
	}

	// Create server.
	server := wsApi.NewServer(nodeConns, localAddress)

	stored := make(chan error)

	go func() {
		log.Println(server.Listen())

		stopped := make(chan error)

		// Emit stop callback.
		go func() {
			stopped <- node.System(creds).Native().EmitEvent(bitnode.LifecycleStop)
		}()

		// Create store.
		st := store.NewStore("store")

		// Store node.
		if err := nodeConns.Store(st); err != nil {
			stored <- err
			return
		}

		// Write node store.
		if err := st.Write("."); err != nil {
			log.Println(err)
			stored <- err
			return
		}

		_ = node.System(creds).Native().EmitEvent(bitnode.LifecycleStop)

		err := <-stopped
		stored <- err
	}()

	log.Printf("Listening on %s...", server.Address())

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	<-cancelChan

	log.Println("Stopping...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}

	if err := <-stored; err != nil {
		log.Printf("Error storing node: %v", err)
	}
}

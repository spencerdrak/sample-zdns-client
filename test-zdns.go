package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	"github.com/liip/sheriff"
	"github.com/zmap/zdns/pkg/zdns"
)

func main() {
	// You'll likely want to replace this with some other type of datastructure in a real application, but for a sample, this works.
	hosts := [1]string{"censys.io"}

	rawAnswers := make([]zdns.Response, 0)

	rawModule := zdns.RawModule{}
	conn, localAddr, err := rawModule.NewReusableUDPConn(nil)

	if err != nil {
		panic(err)
	}

	rawClient := rawModule.NewLookupClient()

	rawOptions := zdns.ClientOptions{
		ReuseSockets:          false,
		IsTraced:              true,
		Verbosity:             3,
		TCPOnly:               false,
		UDPOnly:               false,
		NsResolution:          false,
		LocalAddr:             localAddr,
		Conn:                  &conn,
		Nameserver:            "1.1.1.1:53",
		ModuleOptions:         map[string]string{},
		IsInternallyRecursive: false,
		IterativeOptions:      zdns.IterativeOptions{},
	}

	err = rawClient.Initialize(&rawOptions)

	if err != nil {
		panic(err)
	}

	for _, host := range hosts {
		q := zdns.Question{
			Type:    1,
			Class:   1,
			Name:    host,
			Id:      uuid.New(),
			Timeout: 15,
		}
		// We could use goroutines to do this, but in this simple example, this is fine.
		// Everything on the client is thread-safe

		resp, err := rawClient.DoLookup(q)
		if err != nil {
			panic(err)
		}

		v, _ := version.NewVersion("0.0.0")
		o := &sheriff.Options{
			Groups:     []string{"normal"},
			ApiVersion: v,
		}
		data, err := sheriff.Marshal(o, resp)
		jsonRes, err := json.Marshal(data)
		if err != nil {
			log.Fatal("Unable to marshal JSON result", err)
		}
		fmt.Println(string(jsonRes))

		rawAnswers = append(rawAnswers, resp)
	}
}

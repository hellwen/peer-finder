/*
Copyright 2014 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// A small utility program to lookup hostnames of endpoints in a service.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
	"strings"
	set "peer-finder/myset"
)

const (
    defaultDnsSuffix = "svc.zeusis.com"
    pollPeriod = time.Second * 1
)

var (
	onChange  = flag.String("on-change", "", "Script to run on change, must accept a new line separated list of peers via stdin.")
	onStart   = flag.String("on-start", "", "Script to run on start, must accept a new line separated list of peers via stdin.")
	svc       = flag.String("service", "", "Governing service responsible for the DNS records of the domain this pod is in.")
	namespace = flag.String("ns", "", "The namespace this pod is running in. If unspecified, the POD_NAMESPACE env var is used.")
	dnsSuffix = flag.String("dns-suffix", "", "The dns suffix this pod is running in. If unspecified, the 'svc.cluster.local' env var is used.")
)

func lookup(svcName string) (*set.Set, error) {
	endpoints := set.New()
	_, srvRecords, err := net.LookupSRV("", "", svcName)
	if err != nil {
		return endpoints, err
	}
	for _, srvRecord := range srvRecords {
		// The SRV records ends in a "." for the root domain
		ep := fmt.Sprintf("%v", srvRecord.Target[:len(srvRecord.Target)-1])
        ips, _ := net.LookupIP(ep)
        endpoints.Add(ep + "," + ips[0].String())
	}
	return endpoints, nil
}

func shellOut(sendStdin, script string) {
	log.Printf("execing: %v with stdin: %v", script, sendStdin)
	// TODO: Switch to sending stdin from go

	out, err := exec.Command("bash", "-c", fmt.Sprintf("echo -e '%v' | %v", sendStdin, script)).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute %v: %v, err: %v", script, string(out), err)
	}
	log.Print(string(out))
}

func main() {
	flag.Parse()

	ns := *namespace
	if ns == "" {
		ns = os.Getenv("POD_NAMESPACE")
	}
	suffix := *dnsSuffix
	if suffix == "" {
        	suffix = defaultDnsSuffix
	}
	if *svc == "" || suffix == "" || ns == "" || (*onChange == "" && *onStart == "") {
		log.Fatalf("Incomplete args, see the --help.")
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %s", err)
	}

	myName := strings.Join([]string{hostname, *svc, ns, suffix}, ".")
	ips, _ := net.LookupIP(myName)
	myName = myName + "," + ips[0].String()
	script := *onStart
	if script == "" {
		script = *onChange
		log.Printf("No on-start supplied, on-change %v will be applied on start.", script)
	}
	for newPeers, peers := set.New(), set.New(); script != ""; time.Sleep(pollPeriod) {

		newPeers, err = lookup(*svc)
		if err != nil {
			log.Printf("%v", err)
			continue
		}

		if newPeers.Equal(peers) || !newPeers.Has(myName) {
			continue
        }

        peerList := newPeers.SortList()
        log.Printf("Peer list updated\niam %v\nwas %v\nnow %v", myName, peers.SortList(), peerList)
        shellOut(strings.Join(peerList, "\n"), script)
		peers = newPeers
		script = *onChange
	}
	// TODO: Exit if there's no on-change?
	log.Printf("Peer finder exiting")
}

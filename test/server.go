/*
Copyright 2020 The cert-manager Authors.

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

package test

import (
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	"github.com/miekg/dns"
)

const (
	defaultTTL = 1
)

var requestCount = map[string]int{}

type Handler struct {
	Log logr.Logger

	TxtRecords map[string][][]string
	Zones      []string
	tsigZone   string
	lock       sync.Mutex
}

// serveDNS implements github.com/miekg/dns.Handler
func (b *Handler) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	b.lock.Lock()
	defer b.lock.Unlock()
	log := b.Log.WithName("serveDNS")

	m := new(dns.Msg)
	m.SetReply(req)
	defer w.WriteMsg(m)

	reqName := req.Question[0].Name

	fmt.Printf("Received request for %s\n", reqName)

	if requestCount[reqName] < len(b.TxtRecords[reqName]) {
		for _, record := range b.TxtRecords[reqName][requestCount[reqName]] {
			txtRR, _ := dns.NewRR(fmt.Sprintf("%s %d IN TXT %s", reqName, defaultTTL, record))
			m.Answer = append(m.Answer, txtRR)
		}
		requestCount[reqName]++
	}

	for _, rr := range m.Answer {
		log.Info("responding", "response", rr.String())
	}
}

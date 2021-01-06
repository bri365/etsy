// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	v3 "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/report"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
	"gopkg.in/cheggaaa/pb.v1"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Benchmark put",

	Run: putFunc,
}

var (
	keySize int
	valSize int

	putTotal int
	putRate  int

	keySpaceSize int
	seqKeys      bool

	compactInterval   time.Duration
	compactIndexDelta int64

	checkHashkv bool
)

func init() {
	RootCmd.AddCommand(putCmd)
	putCmd.Flags().IntVar(&keySize, "key-size", 256, "Key size of put request")
	putCmd.Flags().IntVar(&valSize, "val-size", 1024, "Value size of put request")
	putCmd.Flags().IntVar(&putRate, "rate", 0, "Maximum puts per second (0 is no limit)")

	putCmd.Flags().IntVar(&putTotal, "total", 10000, "Total number of put requests")
	putCmd.Flags().IntVar(&keySpaceSize, "key-space-size", 10000, "Maximum possible keys")
	putCmd.Flags().BoolVar(&seqKeys, "sequential-keys", false, "Use sequential keys")
	putCmd.Flags().DurationVar(&compactInterval, "compact-interval", 0, `Interval to compact database (do not duplicate this with etcd's 'auto-compaction-retention' flag) (e.g. --compact-interval=5m compacts every 5-minute)`)
	putCmd.Flags().Int64Var(&compactIndexDelta, "compact-index-delta", 1000, "Delta between current revision and compact revision (e.g. current revision 10000, compact at 9000)")
	putCmd.Flags().BoolVar(&checkHashkv, "check-hashkv", false, "'true' to check hashkv")
}

func putFunc(cmd *cobra.Command, args []string) {
	if keySize < 8 {
		fmt.Fprintf(os.Stderr, "expected --key-size of 8 or more, got (%v)", keySize)
		os.Exit(1)
	}

	if keySpaceSize <= 0 {
		fmt.Fprintf(os.Stderr, "expected positive --key-space-size, got (%v)", keySpaceSize)
		os.Exit(1)
	}

	// Save keys to a file for use in other tests
	name := fmt.Sprintf("keys/keys-%d-%d-%d", time.Now().Unix(), keySize, keySpaceSize)
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var key []byte
	var keys [][]byte
	for i := 0; i < keySpaceSize; i++ {
		data := make([]byte, keySize)
		rand.Read(data)
		keys = append(keys, data)
		bytesWritten, err := f.Write(data)
		if err != nil || bytesWritten != keySize {
			panic(err)
		}
	}

	if seqKeys {
		// sequential keys start at zero for now
		keys[0] = make([]byte, keySize)
	} else {
		k := keys[0]
		// zero out the final uint64 bytes
		binary.LittleEndian.PutUint64(k[keySize-8:keySize], 0)
	}
	// fmt.Printf("%v\n", keys[0])

	value := make([]byte, valueSize)

	requests := make(chan v3.Op, totalClients)
	if putRate == 0 {
		putRate = math.MaxInt32
	}
	limit := rate.NewLimiter(rate.Limit(putRate), 1)
	clients := mustCreateClients(totalClients, totalConns)

	bar = pb.New(putTotal)
	bar.Format("... !")
	bar.Start()

	r := newReport()
	for i := range clients {
		wg.Add(1)
		go func(c *v3.Client) {
			defer wg.Done()
			for op := range requests {
				limit.Wait(context.Background())

				st := time.Now()
				_, err := c.Do(context.Background(), op)
				r.Results() <- report.Result{Err: err, Start: st, End: time.Now()}
				bar.Increment()
			}
		}(clients[i])
	}

	go func() {
		for i := 0; i < putTotal; i++ {
			if seqKeys {
				key = keys[0]
				binary.LittleEndian.PutUint64(key[keySize-8:keySize], uint64(i))
			} else {
				key = keys[i%keySpaceSize]
			}
			// fmt.Printf("\n%v\n", binary.LittleEndian.Uint64(key[keySize-8:keySize]))
			requests <- v3.OpPut(string(key), string(value))
		}
		close(requests)
	}()

	if compactInterval > 0 {
		go func() {
			for {
				time.Sleep(compactInterval)
				compactKV(clients)
			}
		}()
	}

	rc := r.Run()
	wg.Wait()
	close(r.Results())
	bar.Finish()
	fmt.Println(<-rc)

	if checkHashkv {
		hashKV(cmd, clients)
	}
}

func compactKV(clients []*v3.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := clients[0].KV.Get(ctx, "foo")
	cancel()
	if err != nil {
		panic(err)
	}
	revToCompact := max(0, resp.Header.Revision-compactIndexDelta)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = clients[0].KV.Compact(ctx, revToCompact)
	cancel()
	if err != nil {
		panic(err)
	}
}

func max(n1, n2 int64) int64 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func hashKV(cmd *cobra.Command, clients []*v3.Client) {
	eps, err := cmd.Flags().GetStringSlice("endpoints")
	if err != nil {
		panic(err)
	}
	for i, ip := range eps {
		eps[i] = strings.TrimSpace(ip)
	}
	fmt.Println(eps)
	host := eps[0]

	st := time.Now()
	clients[0].HashKV(context.Background(), eps[0], 0)
	rh, eh := clients[0].HashKV(context.Background(), host, 0)
	if eh != nil {
		fmt.Fprintf(os.Stderr, "Failed to get the hashkv of endpoint %s (%v)\n", host, eh)
		panic(err)
	}
	rt, es := clients[0].Status(context.Background(), host)
	if es != nil {
		fmt.Fprintf(os.Stderr, "Failed to get the status of endpoint %s (%v)\n", host, es)
		panic(err)
	}

	rs := "HashKV Summary:\n"
	rs += fmt.Sprintf("\tHashKV: %d\n", rh.Hash)
	rs += fmt.Sprintf("\tEndpoint: %s\n", host)
	rs += fmt.Sprintf("\tTime taken to get hashkv: %v\n", time.Since(st))
	rs += fmt.Sprintf("\tDB size: %s", humanize.Bytes(uint64(rt.DbSize)))
	fmt.Println(rs)
}

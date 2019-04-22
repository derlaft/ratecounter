# ratecounter

## [ratecounter/counter](/counter/)

* Counts increments with given precision `Accuracy` (as `time.Interval`).
* Increments older than `WindowSize` (`time.Interval`) are dropped.
* Fast `Incr()`, `Count()`, `Save()`, `Load()`.
* Low and somewhat constant memory consumption.
* Saves data in JSON format.
* Only standard library is used (except for testing).
* tl;dr - fast enough implementation.

```
  Benchmarks of basic functionality

  Ran 10 samples:
  1kk calls to insert:
    Fastest Time: 0.136s
    Slowest Time: 0.146s
    Average Time: 0.137s ± 0.003s
  1kk calls to count:
    Fastest Time: 0.144s
    Slowest Time: 0.162s
    Average Time: 0.147s ± 0.005s
    
  Benchmarks of i/o
    
  Ran 100 samples:
  Save performance:
    Fastest Time: 0.000s
    Slowest Time: 0.000s
    Average Time: 0.000s ± 0.000s
  Load performance:
    Fastest Time: 0.000s
    Slowest Time: 0.000s
    Average Time: 0.000s ± 0.000s
```

## [ratecounter/nanocounter](/nanocounter/)

* Counts increments with nanosecond precision.
* Increments older than WindowSize are dropped.
* Little slower `Incr()`, `Count()`.
* Much slower `Save()` and `Load()` (a lot of data to save).
* Higher memory consumption, relative to number of increments.
* Saves data in binary format.
* Only standard library is used (except for testing).
* **tl;dr: May be usefull somewhere else** (storing some additional data with timestamps?), **but not in this sample HTTP server (works only a little slower, but requires much higher IO usage).**

```
  Benchmarks of basic functionality

  Ran 10 samples:
  Lots of calls to insert:
    Fastest Time: 0.131s
    Slowest Time: 0.170s
    Average Time: 0.143s ± 0.011s
  Lots of calls to count:
    Fastest Time: 0.201s
    Slowest Time: 0.211s
    Average Time: 0.204s ± 0.003s

Benchmarks of i/o

  Ran 100 samples:
  Lots of calls to save:
    Fastest Time: 0.028s
    Slowest Time: 0.044s
    Average Time: 0.033s ± 0.003s
  Lots of calls to load:
    Fastest Time: 0.110s
    Slowest Time: 0.148s
    Average Time: 0.117s ± 0.006s
```

## [ratecounter/external](/external/)

* Really wanted to compare my results to some existing code.
* Contains wrapper to compare performance with an [existing open-source implementation](https://github.com/paulbellamy/ratecounter).
* Constant-time and the fastest `Count()`
* Fastest `Incr()`.
* `Save()` and `Load()` might be tricky to implement, but still possible.
* Uses a timer in a separate routine to handle which [may be a problem in some rare conditions](https://github.com/paulbellamy/ratecounter/issues/14).
* No dynamic allocations.
* tl;dr: Superior library, ~~hire the author instead~~ should probably use it in production environment instead %)

```
  Benchmarks of basic functionality

  Ran 10 samples:
  1kk calls to insert:
    Fastest Time: 0.024s
    Slowest Time: 0.028s
    Average Time: 0.025s ± 0.002s
  1kk calls to count:
    Fastest Time: 0.002s
    Slowest Time: 0.004s
    Average Time: 0.003s ± 0.000s
```

# Sample test tool - HTTP server

* Installation - `go get github.com/derlaft/ratecounter/cmd/rateserver`.
* Run - `rateserver`.
* Uses the first implementation (`ratecounter/counter`).
* Listens on `127.0.0.1:8081` (`listenAddr` in [main.go](/cmd/rateserver/main.go))
* Counts client requests (`windowSize=60s`, `accuracy=200ms` in [main.go](/cmd/rateserver/main.go)).
* Saves data to disk on exit signal (`filename=state.rtt` in [main.go](/cmd/rateserver/main.go)).
* Saves data to disk on a timer (`checkpointInterval=500ms` in [aux.go](/cmd/rateserver/aux.go)).
* Only standard library is used (except for testing).
* Request example:

```bash
% curl http://localhost:8081; echo
00000000000000000017
```

## Sample benchmark

```bash
% ab -n 1000000 -c 100 "http://localhost:8081/"                                                                        :(

...

Server Software:        
Server Hostname:        localhost
Server Port:            8081

Document Path:          /
Document Length:        20 bytes

Concurrency Level:      100
Time taken for tests:   62.540 seconds
Complete requests:      1000000
Failed requests:        0
Total transferred:      137000000 bytes
HTML transferred:       20000000 bytes
Requests per second:    15989.75 [#/sec] (mean)
Time per request:       6.254 [ms] (mean)
Time per request:       0.063 [ms] (mean, across all concurrent requests)
Transfer rate:          2139.25 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    3   0.5      3      11
Processing:     0    4   1.3      3      75
Waiting:        0    3   1.3      2      75
Total:          0    6   1.2      6      75

Percentage of the requests served within a certain time (ms)
  50%      6
  66%      6
  75%      7
  80%      7
  90%      7
  95%      8
  98%      9
  99%     10
 100%     75 (longest request)
```

## Misc

* No channels or coroutines [were used](cmd/rateserver/main.go#L42) in the counter implementation.
* No `sync/atomic` was used. This is not [a special, low-level application](https://golang.org/pkg/sync/atomic/). It also a common source of annoying bugs (okay, at least for me and people I worked with).
* Full and correct gracefull termination of all modules is important and should be implemented.

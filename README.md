# ratecounter

## [ratecounter/counter](/counter/)

* Counts increments with given precision `Accuracy` (as `time.Interval`).
* Increments older than `WindowSize` (`time.Interval`) are dropped.
* Fast `Incr()`, `Count()`, `Save()`, `Load()`.
* Low and somewhat constant memory consumption.
* Saves data in JSON format.
* Only standard library is used (except for testing).

## [ratecounter/nanocounter](/nanocounter/)

* Counts increments with nanosecond precision.
* Increments older than WindowSize are dropped.
* Little slower `Incr()`, `Count()`.
* Much slower `Save()` and `Load()` (a lot of data to save).
* Higher memory consumption, relative to number of increments.
* Saves data in binary format.
* Only standard library is used (except for testing).
* **May be usefull somewhere else, but not in this sample HTTP server (works only a little slower, but requires much higher IO usage).**

# Sample test tool - HTTP server

* Installation - `go get github.com/derlaft/ratecounter/cmd/rateserver`.
* Run - `rateserver`.
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
This is ApacheBench, Version 2.3 <$Revision: 1843412 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 100000 requests
Completed 200000 requests
Completed 300000 requests
Completed 400000 requests
Completed 500000 requests
Completed 600000 requests
Completed 700000 requests
Completed 800000 requests
Completed 900000 requests
Completed 1000000 requests
Finished 1000000 requests


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

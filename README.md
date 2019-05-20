# ratecounter

# Sample test tool - HTTP server

* Installation - `go get github.com/derlaft/ratecounter/cmd/rateserver`.
* Run - `rateserver`.
* Uses the first implementation (`ratecounter/counter`).
    * The implementation was changed so that it does not directly saves data to disk, but rather returns raw values.
    * Factories were added so the new code module is tested clearly.
* Listens on `127.0.0.1:8081` (`listenAddr` in [main.go](/cmd/rateserver/main.go))
* Counts total client requests (`windowSize=20s`, `accuracy=200ms` in [main.go](/cmd/rateserver/main.go)).
* Rejects reqeusts from IP if there were more of them in a given amount of time (`maxRequests=15` in [main.go](/cmd/rateserver/main.go)).
* Supports `X-Real-IP` HTTP header (see example).
    * Will probably have some weird problems with address detection with IPv6.
* Saves data to disk on exit signal (`filename=state.rtt` in [main.go](/cmd/rateserver/main.go)).
* Saves data to disk on a timer (`checkpointInterval=500ms` in [aux.go](/cmd/rateserver/aux.go)).
* Only standard library is used (except for testing).
* Request example:

```bash
% for i in $(seq 1 20); do curl -H 'X-Real-IP: 11.8.8.9' http://localhost:8081; done
00001
00002
00003
00004
00005
00006
00007
00008
00009
00010
00011
00012
00013
00014
00015
Error while processing your request: Too many requests from your IP, go away
Error while processing your request: Too many requests from your IP, go away
Error while processing your request: Too many requests from your IP, go away
Error while processing your request: Too many requests from your IP, go away
Error while processing your request: Too many requests from your IP, go away
% for i in $(seq 1 20); do curl -H 'X-Real-IP: 11.8.8.10' http://localhost:8081; done
00016
00017
00018
00019
00020
00021
00022
00023
00024
00025
00026
00027
00028
00029
00030
Error while processing your request: Too many requests from your IP, go away
Error while processing your request: Too many requests from your IP, go away
Error while processing your request: Too many requests from your IP, go away
Error while processing your request: Too many requests from your IP, go away
Error while processing your request: Too many requests from your IP, go away
```

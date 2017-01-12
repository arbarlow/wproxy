# Weave Proxy

The proxy is ritten in Go and only uses the stdlib. You should be able to...

```
get get github.com/arbarlow/wproxy
wproxy -h
```

To accept connections on 8000 and foward to 8001 run

```
wproxy -listen :8000 -server :8001
```

### Thoughts

The proxy uses a tee reader to funnel the tcp connection to a byte scanner, I figured there was no need to write a "parser" since I was only actually interested in new lines and the first char of each line.

A ticker created new array elements every second and the stats package attomically updates the last elements in these arrays to provide stats.

There are no tests, usually I would TDD but given the client and server binaries it was a fairly testable run and repeat pattern

Possible extra stats, Failure percentage?

### Issues 

Closing connection creates error, probably due to server response of closed connection and timing, I moved on rather than tracking it down for now

Timing stats array could get very large, really I think the ticker should calculate a running average and remove old array values

Any other thoughts and discussion are welcome, it was rushed.

# Xmux

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/rs/xmux) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/xmux/master/LICENSE) [![Build Status](https://travis-ci.org/rs/xmux.svg?branch=master)](https://travis-ci.org/rs/xmux) [![Coverage](http://gocover.io/_badge/github.com/rs/xmux)](http://gocover.io/github.com/rs/xmux)

Xmux is a lightweight high performance HTTP request muxer on top [xhandler](https://github.com/rs/xhandler). Xmux gets its speed from the fork of the amazing [httprouter](https://github.com/julienschmidt/httprouter). Route parameters are stored in `context` instead of being passed as an additional parameter.

In contrast to the [default mux](http://golang.org/pkg/net/http/#ServeMux) of Go's `net/http` package, this muxer supports variables in the routing pattern and matches against the request method. It also scales better.

The muxer is optimized for high performance and a small memory footprint. It scales well even with very long paths and a large number of routes. A compressing dynamic trie (radix tree) structure is used for efficient matching.

## Features

**Only explicit matches:** With other muxers, like [http.ServeMux](http://golang.org/pkg/net/http/#ServeMux), a requested URL path could match multiple patterns. Therefore they have some awkward pattern priority rules, like *longest match* or *first registered, first matched*. By design of this router, a request can only match exactly one or no route. As a result, there are also no unintended matches, which makes it great for SEO and improves the user experience.

**Stop caring about trailing slashes:** Choose the URL style you like, the muxer automatically redirects the client if a trailing slash is missing or if there is one extra. Of course it only does so, if the new path has a handler. If you don't like it, you can [turn off this behavior](http://godoc.org/github.com/rs/xmux#Mux.RedirectTrailingSlash).

**Path auto-correction:** Besides detecting the missing or additional trailing slash at no extra cost, the muxer can also fix wrong cases and remove superfluous path elements (like `../` or `//`). Is [CAPTAIN CAPS LOCK](http://www.urbandictionary.com/define.php?term=Captain+Caps+Lock) one of your users? Xmux can help him by making a case-insensitive look-up and redirecting him to the correct URL.

**Parameters in your routing pattern:** Stop parsing the requested URL path, just give the path segment a name and the router delivers the dynamic value to you. Because of the design of the router, path parameters are very cheap.

**RouteGroups:** A way to create [groups of routes](http://godoc.org/github.com/rs/xmux#Mux.NewGroup) without incurring any per-request overhead.

**Zero Garbage:** The matching and dispatching process generates zero bytes of garbage. In fact, the only heap allocations that are made, is by building the slice of the key-value pairs for path parameters and the `context` instance to store them in the context. If the request path contains no parameters, not a single heap allocation is necessary.

**No more server crashes:** You can set a [Panic handler](http://godoc.org/github.com/rs/xmux#Mux.PanicHandler) to deal with panics occurring during handling a HTTP request. The router then recovers and lets the `PanicHandler` log what happened and deliver a nice error page.

Of course you can also set **custom [NotFound](http://godoc.org/github.com/rs/xmux#Mux.NotFound) and  [MethodNotAllowed](http://godoc.org/github.com/rs/xmux#Mux.MethodNotAllowed) handlers**.

## Usage

This is just a quick introduction, view the [GoDoc](http://godoc.org/github.com/rs/xmux) for details.

Let's start with a trivial example:
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"context"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %s!\n", xmux.Param(ctx, "name"))
}

func main() {
	mux := xmux.New()
	mux.GET("/", xhandler.HandlerFuncC(Index))
	mux.GET("/hello/:name", xhandler.HandlerFuncC(Hello))

	log.Fatal(http.ListenAndServe(":8080", xhandler.New(context.Background(), mux)))
}
```

You may also chain middleware using `xhandler.Chain`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

func main() {
	c := xhandler.Chain{}

	// Append a context-aware middleware handler
	c.UseC(xhandler.CloseHandler)

	// Another context-aware middleware handler
	c.UseC(xhandler.TimeoutHandler(2 * time.Second))

	mux := xmux.New()

	// Use c.Handler to terminate the chain with your final handler
	mux.GET("/welcome/:name", xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome %s!", xmux.Param(ctx, "name"))
	}))

	if err := http.ListenAndServe(":8080", c.Handler(mux)); err != nil {
		log.Fatal(err)
	}
}
```

### Named parameters

As you can see, `:name` is a *named parameter*. The values are accessible via `xmux.Params(ctx)`, which returns `xmux.ParamHolder`.
You can get the value of a parameter by its name using `Get(name)` method:

```go
user := xmux.Params(ctx).Get("user")
```

or using `xmux.Param(ctx, name)` shortcut:

```go
user := xmux.Param(ctx, "user")
```

Named parameters only match a single path segment:

```
Pattern: /user/:user

 /user/gordon              match
 /user/you                 match
 /user/gordon/profile      no match
 /user/                    no match
```

**Note:** Since this muxer has only explicit matches, you can not register static routes and parameters for the same path segment. For example you can not register the patterns `/user/new` and `/user/:user` for the same request method at the same time. The routing of different request methods is independent from each other.

### Catch-All parameters

The second type are *catch-all* parameters and have the form `*name`. Like the name suggests, they match everything. Therefore they must always be at the **end** of the pattern:

```
Pattern: /src/*filepath

 /src/                     match
 /src/somefile.go          match
 /src/subdir/somefile.go   match
```

## Benchmarks

Thanks to [Julien Schmidt](https://github.com/julienschmidt) excellent [HTTP routing benchmark](https://github.com/julienschmidt/go-http-routing-benchmark), we can see that xhandler's muxer is pretty close to `httprouter` as it is a fork of it. The small overhead is due to the `context` allocation used to store route parameters. It still outperform other routers, thanks to amazing `httprouter`'s radix tree based matcher.

```
BenchmarkXhandler_APIStatic-8   	50000000	        39.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkChi_APIStatic-8        	 3000000	       439 ns/op	     144 B/op	       5 allocs/op
BenchmarkGoji_APIStatic-8       	 5000000	       272 ns/op	       0 B/op	       0 allocs/op
BenchmarkHTTPRouter_APIStatic-8 	50000000	        37.3 ns/op	       0 B/op	       0 allocs/op

BenchmarkXhandler_APIParam-8    	 5000000	       328 ns/op	     160 B/op	       4 allocs/op
BenchmarkChi_APIParam-8         	 2000000	       675 ns/op	     432 B/op	       6 allocs/op
BenchmarkGoji_APIParam-8        	 2000000	       692 ns/op	     336 B/op	       2 allocs/op
BenchmarkHTTPRouter_APIParam-8  	10000000	       166 ns/op	      64 B/op	       1 allocs/op

BenchmarkXhandler_API2Params-8  	 5000000	       362 ns/op	     160 B/op	       4 allocs/op
BenchmarkChi_API2Params-8       	 2000000	       814 ns/op	     432 B/op	       6 allocs/op
BenchmarkGoji_API2Params-8      	 2000000	       680 ns/op	     336 B/op	       2 allocs/op
BenchmarkHTTPRouter_API2Params-8	10000000	       183 ns/op	      64 B/op	       1 allocs/op

BenchmarkXhandler_APIAll-8      	  200000	      6473 ns/op	    2176 B/op	      64 allocs/op
BenchmarkChi_APIAll-8           	  100000	     17261 ns/op	    8352 B/op	     146 allocs/op
BenchmarkGoji_APIAll-8          	  100000	     15052 ns/op	    5377 B/op	      32 allocs/op
BenchmarkHTTPRouter_APIAll-8    	  500000	      3716 ns/op	     640 B/op	      16 allocs/op

BenchmarkXhandler_Param1-8      	 5000000	       271 ns/op	     128 B/op	       4 allocs/op
BenchmarkChi_Param1-8           	 2000000	       620 ns/op	     432 B/op	       6 allocs/op
BenchmarkGoji_Param1-8          	 3000000	       522 ns/op	     336 B/op	       2 allocs/op
BenchmarkHTTPRouter_Param1-8    	20000000	       112 ns/op	      32 B/op	       1 allocs/op

BenchmarkXhandler_Param5-8      	 3000000	       414 ns/op	     256 B/op	       4 allocs/op
BenchmarkChi_Param5-8           	 1000000	      1204 ns/op	     432 B/op	       6 allocs/op
BenchmarkGoji_Param5-8          	 2000000	       847 ns/op	     336 B/op	       2 allocs/op
BenchmarkHTTPRouter_Param5-8    	 5000000	       247 ns/op	     160 B/op	       1 allocs/op

BenchmarkXhandler_Param20-8     	 2000000	       747 ns/op	     736 B/op	       4 allocs/op
BenchmarkChi_Param20-8          	 2000000	       746 ns/op	     736 B/op	       4 allocs/op
BenchmarkGoji_Param20-8         	  500000	      2439 ns/op	    1247 B/op	       2 allocs/op
BenchmarkHTTPRouter_Param20-8   	 3000000	       585 ns/op	     640 B/op	       1 allocs/op

BenchmarkXhandler_ParamWrite-8  	 5000000	       404 ns/op	     144 B/op	       5 allocs/op
BenchmarkChi_ParamWrite-8       	 3000000	       407 ns/op	     144 B/op	       5 allocs/op
BenchmarkGoji_ParamWrite-8      	 2000000	       594 ns/op	     336 B/op	       2 allocs/op
BenchmarkHTTPRouter_ParamWrite-8	10000000	       166 ns/op	      32 B/op	       1 allocs/op
```

You can run this benchmark by executing the following commands at the root of `xmux` repository:

```go
go get -t ./bench/routers
go test ./bench/routers -bench .
```

## How does it work?

The muxer relies on a tree structure which makes heavy use of *common prefixes*, it is basically a *compact* [*prefix tree*](http://en.wikipedia.org/wiki/Trie) (or just [*Radix tree*](http://en.wikipedia.org/wiki/Radix_tree)). Nodes with a common prefix also share a common parent. Here is a short example what the routing tree for the `GET` request method could look like:

```
Priority   Path             Handle
9          \                *<1>
3          ├s               nil
2          |├earch\         *<2>
1          |└upport\        *<3>
2          ├blog\           *<4>
1          |    └:post      nil
1          |         └\     *<5>
2          ├about-us\       *<6>
1          |        └team\  *<7>
1          └contact\        *<8>
```

Every `*<num>` represents the memory address of a handler function (a pointer). If you follow a path trough the tree from the root to the leaf, you get the complete route path, e.g `\blog\:post\`, where `:post` is just a placeholder ([*parameter*](#named-parameters)) for an actual post name. Unlike hash-maps, a tree structure also allows us to use dynamic parts like the `:post` parameter, since we actually match against the routing patterns instead of just comparing hashes. [As benchmarks show](https://github.com/julienschmidt/go-http-routing-benchmark), this works very well and efficient.

Since URL paths have a hierarchical structure and make use only of a limited set of characters (byte values), it is very likely that there are a lot of common prefixes. This allows us to easily reduce the routing into ever smaller problems. Moreover the router manages a separate tree for every request method. For one thing it is more space efficient than holding a method->handle map in every single node, for another thing is also allows us to greatly reduce the routing problem before even starting the look-up in the prefix-tree.

For even better scalability, the child nodes on each tree level are ordered by priority, where the priority is just the number of handles registered in sub nodes (children, grandchildren, and so on..). This helps in two ways:

1. Nodes which are part of the most routing paths are evaluated first. This helps to make as much routes as possible to be reachable as fast as possible.
2. It is some sort of cost compensation. The longest reachable path (highest cost) can always be evaluated first. The following scheme visualizes the tree structure. Nodes are evaluated from top to bottom and from left to right.

```
├------------
├---------
├-----
├----
├--
├--
└-
```

## Why doesn't this work with http.Handler?

**It does!** The router itself implements the http.Handler interface. Moreover the router provides convenient [adapters for http.Handler](http://godoc.org/github.com/rs/xmux#Mux.Handle)s and [http.HandlerFunc](http://godoc.org/github.com/rs/xmux#Mux.HandleFunc)s which allows them to be used as a [xhandler.HandlerC](http://godoc.org/github.com/rs/xhandler#HandlerC) when registering a route. The only disadvantage is, that no `context` and thus no parameter values can be retrieved when a `http.Handler` or `http.HandlerFunc` is used.

## Where can I find Middleware *X*?

This package just provides a very efficient request muxer with a few extra features. The muxer is just a [xhandler.HandlerC](https://godoc.org/github.com/rs/xhandler#HandlerC), you can chain any `http.Handler` or `xhandler.HandlerC` compatible middleware before the router, for example the [Gorilla handlers](http://www.gorillatoolkit.org/pkg/handlers). Or you could [just write your own](http://justinas.org/writing-http-middleware-in-go/), it's very easy!

### Multi-domain / Sub-domains

Here is a quick example: Does your server serve multiple domains / hosts? You want to use sub-domains? Define a router per host!

```go
// We need an object that implements the xhandler.HandlerC interface.
// Therefore we need a type for which we implement the ServeHTTP method.
// We just use a map here, in which we map host names (with port) to xhandler.HandlerC
type HostSwitch map[string]xhandler.HandlerC

// Implement the ServerHTTP method on our new type
func (hs HostSwitch) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Check if a xhandler.HandlerC is registered for the given host.
	// If yes, use it to handle the request.
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTPC(ctx, w, r)
	} else {
		// Handle host names for wich no handler is registered
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

func main() {
	c := xhandler.Chain{}

	// Initialize a muxer as usual
	mux := xmux.New()
	mux.GET("/", Index)
	mux.GET("/hello/:name", Hello)

	// Make a new HostSwitch and insert the muxer (our http handler)
	// for example.com and port 12345
	hs := make(HostSwitch)
	hs["example.com:12345"] = mux

	// Use the HostSwitch to listen and serve on port 12345
	if err := http.ListenAndServe(":12345", c.Handler(hs)); err != nil {
		log.Fatal(err)
	}
}
```

### Basic Authentication
Another quick example: Basic Authentication (RFC 2617) for handles:

```go
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

func BasicAuth(user, pass []byte, next xhandler.HandlerC) xhandler.HandlerC {
	return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		const basicAuthPrefix string = "Basic "

		// Get the Basic Authentication credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 &&
					bytes.Equal(pair[0], user) &&
					bytes.Equal(pair[1], pass) {

					// Delegate request to the next handler
					next.ServeHTTPC(ctx, w, r)
					return
				}
			}
		}

		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Not protected!\n")
}

func Protected(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Protected!\n")
}

func main() {
	user := []byte("gordon")
	pass := []byte("secret!")

	c := xhandler.Chain{}
	mux := xmux.New()
	mux.GET("/", xhandler.HandlerFuncC(Index))
	mux.GET("/protected/", BasicAuth(user, pass, xhandler.HandlerFuncC(Protected)))

	log.Fatal(http.ListenAndServe(":8080", c.Handler(mux)))
}
```

## Licenses

All source code is licensed under the [BSD License](https://raw.github.com/rs/xhandler/master/LICENSE).

Xmux is forked from [httprouter](https://github.com/julienschmidt/httprouter) with [BSD License](https://github.com/julienschmidt/httprouter/blob/master/LICENSE).

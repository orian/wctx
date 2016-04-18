# With ConTeXt
[![GoDoc](http://godoc.org/github.com/orian/wctx?status.png)](http://godoc.org/github.com/orian/wctx)

It's a `github.com/julienschmidt/httprouter` drop in replacement. Basically all methods from the original library are wrapped.
It makes easy to us the `httprouter` on App Engine and outside of App Engine.

The problem with App Engine context is that it requires your app to hardcode the `appengine.NewContext` in 
your handlers. As stated in https://blog.golang.org/context inside Google's infrastructure
they pass Context around. It's pretty useful for puting timeouts and other limits on a request
and passing user permissions.

On App Engine the init code would look:

```go
func SomeHandler(c context.Context, w http.ResponseWriter, r *http.Request) { /* some code */ }

func GaeContext(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

func init() {
	r := New()
	r.ContextFactory = GaeContext
	r.GET("/page/:pageid/:mod", SomeHandler)
	http.Handler("", r)  // r may be replaced by r.R, the underlaying httprouter.
}
```

For more details please the doc of the original repo: [![GoDoc](http://godoc.org/github.com/julienschmidt/httprouter?status.png)](http://godoc.org/github.com/julienschmidt/httprouter)

# Other libraries:
 
 - context aware sessions: [github.com/orian/sessions](https://github.com/orian/sessions), based on Gorilla's session library
 - CSRF protection, [github.com/orian/nosurf](https://github.com/orian/nosurf)
 - log library, wrapper for both Google App Engine and standard log, [github.com/orian/utils/net/log](https://github.com/orian/utils/net/log)
 - a User providing library [coming soon]
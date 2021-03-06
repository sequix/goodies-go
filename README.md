# goodies-go
Reuseable code pieces for golang.

## Good Stuff & Advises

### Code generators

1. Regexp debugger & generator: [regex101.com](https://regex101.com)
2. JSON to go struct: [mholt.github.io/json-to-go](https://mholt.github.io/json-to-go)
3. YAML to go struct: [yaml2go.prasadg.dev](https://yaml2go.prasadg.dev)
4. cURL to go code: [curl.trillworks.com/#go](https://curl.trillworks.com/#go)

### Errors

1. Use %s to print err.
2. Do not use word 'failed', 'unsuccessful' and etc, the error itself says it is a error already.
3. Wrap errors as each layer the error passed on.
4. Each layer's wrapping states only what it is doing.

### Cache

1. [patrick/go-cache](https://github.com/patrickmn/go-cache): programmer-friendly in-memory cache, see also [this](https://sequix.cn/post/2020-01-25-valyala-fastcache/).
2. [VictoriaMetrics/fastcache](https://github.com/VictoriaMetrics/fastcache): fixed-size in-memory cache based on ringbuffer

### CLI

Just [cobra](https://github.com/spf13/cobra), and check generator usage [doc](https://github.com/spf13/cobra/blob/master/cobra/README.md).
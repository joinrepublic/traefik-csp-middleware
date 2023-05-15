# Content Security Policy rewrite for Traefik

This is a fork of [Rewrite Body with compression](https://github.com/packruler/rewrite-body) middleware. Instead of rewriting random sequence of strings, it looks for a specific substring in Content Security Policy. It works similar to [CloudFlare Worker for CSP nonces](https://github.com/moveyourdigital/cloudflare-worker-csp-nonce) and is designed to be compatible.


### Process For Handling Body Content

#### Body Content Requirements

* The header must have `Content-Type` that includes `text`. For example:
  * `text/html`
  * `text/json`
* The header must have `Content-Encoding` header that is supported by this plugin
  * The original plugin supported `Content-Encoding` of `identity` or empty
  * This plugin adds support for `gzip` and `zlib` encoding

#### Processing Paths

* If the either of the previous conditions failes the body is passed on as is and no further processing from this plugin occurs.

* If the `Content-Encoding` is empty or `identity` it is handled in mostly the same manner as the original plugin.

* If the `Content-Encoding` is `gzip` the following process happens:
  * The body content is decompressed by [Go-lang's gzip library](https://pkg.go.dev/compress/gzip)
  * The resulting content is run through the `regex` process created by the original plugin
  * The processed content is then compressed with the same library and returned

## Configuration

### Static

```yaml
pilot:
  token: "xxxx"

experimental:
    plugins:
        rewrite-body-csp:
            moduleName: "github.com/joinrepublic/traefik-csp-middleware"
            version: "v2.0.0"
```

### Dynamic

To configure the `Rewrite Body` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example creates
and uses the `rewritebody` middleware plugin to replace all foo occurences by bar in the HTTP response body.

If you want to apply some limits on the response body, you can chain this middleware plugin with the [Buffering middleware](https://docs.traefik.io/middlewares/buffering/) from Traefik.

```yaml
http:
  routers:
    my-router:
      rule: "Host(`localhost`)"
      middlewares: 
        - "rewrite-foo"
      service: "my-service"

  middlewares:
    rewrite-foo:
      plugin:
        rewrite-body-csp:
          # Keep Last-Modified header returned by the HTTP service.
          # By default, the Last-Modified header is removed.
          lastModified: true

          placeholder: DhcnhD3khTMePgXw

          # logLevel is optional, defaults to Info level.
          # Available logLevels: (Trace: -2, Debug: -1, Info: 0, Warning: 1, Error: 2)
          logLevel: 0

          # monitoring is optional, defaults to below configuration
          # monitoring configuration limits the HTTP queries that are checked for regex replacement.
          monitoring:
            # methods is a string list. Options are standard HTTP Methods. Entries MUST be ALL CAPS
            # For a list of options: https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
            methods:
              - GET
            # types is a string list. Options are HTTP Content Types. Entries should match standard formatting
            # For a list of options: https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
            # Wildcards(*) are not supported!
            types:
              - text/html
  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1"
```


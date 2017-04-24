This is a toy but kinda serious little program to force the "api" entry for a dns zone
to be proxied. Mainly for use as a demonstration of how to do more.

Useful for about 2 more weeks until ansible adds "proxied" to the cloudflare_dns module. ;-)

# Development
Since this is a first and a toy, I feel ok putting some extra details here.

To compile fast, `go build` from inside src/cfpu

To compile for deployment, `go build --ldflags '-s -w '`.

Be sure to update the version number before deploying any changes.

# timeout-tricker
A reverse proxy that tricks your ELB into not timing out your connection.
It only works on requests which  can tolerate spaces prefixed to the body e.g. `json` or `html`
Also, it disables compression for convenience right now. Ultimately it should
include some heuristics or hinting based on the `path` so it can anticipate
the expected response.

# configuration
| Environment | Description |
|-------------|-------------|
| HOST | Upstream host to proxy |
| TIMEOUT | Number of seconds to wait before tricking starts |

# possible TODOs
-	Add some heuristics based on original `path` of the request to anticipate the type of response coming back
-	Buffer or store the response temporarily in case of binary BODY and generate a HTML redirect on-the-fly to this location


# contact / getting help

* andy.lo-a-foe@philips.com

# license
License is MIT

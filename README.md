# timeout-tricker
A reverse proxy that tricks your ELB into not timing out your connection.
It only works on requests that can tolerate spaces prefixed to the body.
Also, it disable compression for convenience right now. Ultimately it should
include some heuristics or hinting based on the `path` so it can anticipate
the expected response.

# configuration
| Environment | Description |
|-------------|-------------|
| HOST | Upstream host to proxy |
| TIMEOUT | Number of seconds to wait before tricking starts |

# Contact / getting help

* andy.lo-a-foe@philips.com

# license
License is MIT

import yaml
import json

schema Server:
    ports: [int]

    check:
        all p in ports {
            0 < p < 65535
        }

server1: Server = yaml.decode("""\
ports:
- 80
- 8080
""")
server2: Server = json.decode("""\
{
    "ports": [80, 8000]
}
""")

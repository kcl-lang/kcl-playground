
from pluto_client import CloudEvent

import json

from pluto_client import KVStore
kvstore = KVStore.build_client("kvstore")
def handle_queue_event(evt: CloudEvent):
    data = json.loads(evt.data)
    print(data)
    kvstore.set(data["name"], data["message"])

_default = handle_queue_event
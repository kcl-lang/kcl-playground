
from pluto_client import HttpRequest

from pluto_client import HttpResponse

from pluto_client import KVStore
kvstore = KVStore.build_client("kvstore")

def store_handler(req: HttpRequest) -> HttpResponse:
    name = req.query["name"] if "name" in req.query else "Anonym"
    if isinstance(name, list):
        name = ",".join(name)
    message = kvstore.get(str(name))
    return HttpResponse(
        status_code=200, body=f"Fetch {name} access message: {message}."
    )

_default = store_handler
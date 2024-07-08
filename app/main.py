import json
import re
import hashlib
from pluto_client import Website, KVStore, Router, HttpRequest, HttpResponse
import kcl_lib.api as kcl_api
import tempfile

color_pattern = re.compile(r"\x1b\[[0-9;]+m")
api = kcl_api.API()
router = Router("router")
store = KVStore("store")
website = Website("./website", "kcl-playground")
website.addEnv("BACKEND_URL", router.url())


def run_code(code: str) -> kcl_api.ExecProgram_Result:
    with tempfile.NamedTemporaryFile(suffix=".k") as temp_file:
        temp_file.write(code.encode())
        temp_file.seek(0)
        args = kcl_api.ExecProgram_Args(k_filename_list=[temp_file.name])
        result = api.exec_program(args)
        return result


def fmt_code(code: str) -> str:
    args = kcl_api.FormatCode_Args(source=code)
    result = api.format_code(args)
    return str(result.formatted, encoding="utf-8")


def compile_handler(req: HttpRequest) -> HttpResponse:
    code = req.body["body"]
    try:
        result = run_code(code)
    except Exception as err:
        return HttpResponse(
            status_code=200,
            body=json.dumps(
                {
                    "errors": color_pattern.sub("", str(err).removeprefix("ERROR:")),
                }
            ),
        )
    if result.err_message:
        return HttpResponse(
            status_code=200,
            body=json.dumps(
                {
                    "errors": color_pattern.sub("", result.err_message),
                }
            ),
        )
    else:
        return HttpResponse(
            status_code=200,
            body=json.dumps(
                {
                    "events": [
                        {
                            "message": result.yaml_result,
                            "kind": "stdout",
                        }
                    ],
                }
            ),
        )


def fmt_handler(req: HttpRequest) -> HttpResponse:
    code = req.body["body"]
    result = fmt_code(code)
    return HttpResponse(status_code=200, body=json.dumps({"body": result}))


def share_handler(req: HttpRequest) -> HttpResponse:
    code = req.body["body"]
    sha1 = hashlib.sha1()
    sha1.update(code.encode("utf-8"))
    id = sha1.hexdigest()
    store.set(id, code)
    return HttpResponse(status_code=200, body=json.dumps({"body": id}))


def query_handler(req: HttpRequest) -> HttpResponse:
    id = req.query["id"]
    return HttpResponse(status_code=200, body=json.dumps({"body": store.get(id)}))


router.post("/-/play/compile", compile_handler)
router.post("/-/play/fmt", fmt_handler)
router.post("/-/play/share", share_handler)
router.get("/-/play/query", query_handler)

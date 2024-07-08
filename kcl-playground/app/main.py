import json
import hashlib
from pluto_client import Website, KVStore, Router, HttpRequest, HttpResponse
import kcl_lib.api as kcl_api

api = kcl_api.API()
website = Website("./website", "kcl-playground")
router = Router("router")
store = KVStore("store")

website.addEnv("BACKEND_URL", router.url())


def run_code(code: str) -> kcl_api.ExecProgram_Result:
    args = kcl_api.ExecProgram_Args(k_filename_list=["test.k"], k_code_list=[code])
    result = api.exec_program(args)
    return result


def fmt_code(code: str) -> str:
    args = kcl_api.FormatCode_Args(source=code)
    result = api.format_code(args)
    return result.formatted


def compile_handler(req: HttpRequest) -> HttpResponse:
    code = req.body["body"]
    result = run_code(code)
    if result.err_message:
        return HttpResponse(
            status_code=200,
            body=json.dumps(
                {
                    "errors": result.err_message,
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
    return HttpResponse(status_code=200, body=json.dumps({"id": id}))


router.post("/-/play/compile", compile_handler)
router.post("/-/play/fmt", fmt_handler)
router.post("/-/play/share", share_handler)

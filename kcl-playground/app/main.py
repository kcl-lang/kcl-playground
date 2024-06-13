import json
from pluto_client import Website, Router, HttpRequest, HttpResponse
import kcl_lib.api as kcl_api

website = Website("./website", "kcl-playground")
router = Router("router")

website.addEnv("BACKEND_URL", router.url())

def run_code(code: str) -> kcl_api.ExecProgram_Result:
    args = kcl_api.ExecProgram_Args(k_filename_list=["test.k"], k_code_list=[code])
    api = kcl_api.API()
    result = api.exec_program(args)
    return result


def compile_handler(req: HttpRequest) -> HttpResponse:
    code = req.body["body"]
    result = run_code(code)
    if result.err_message:
        return HttpResponse(status_code=200, body=json.dumps({
            "errors": result.err_message,
        }))
    else:
        return HttpResponse(status_code=200, body=json.dumps({
            "events": [{
                "message": result.yaml_result,
                "kind": "stdout",
            }],
        }))


router.post("/-/play/compile", compile_handler)

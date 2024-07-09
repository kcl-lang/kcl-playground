import json
import re
import tempfile
from pluto_client import Website, Router, HttpRequest, HttpResponse
import kcl_lib.api as kcl_api

color_pattern = re.compile(r"\x1b\[[0-9;]+m")
api = kcl_api.API()
router = Router("router")
website = Website("./web/dist", "kcl-playground")
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
                    "error": color_pattern.sub("", str(err).removeprefix("ERROR:")),
                }
            ),
        )
    if result.err_message:
        return HttpResponse(
            status_code=200,
            body=json.dumps(
                {
                    "error": color_pattern.sub("", result.err_message),
                }
            ),
        )
    else:
        return HttpResponse(
            status_code=200,
            body=json.dumps({"body": result.yaml_result}),
        )


def fmt_handler(req: HttpRequest) -> HttpResponse:
    code = req.body["body"]
    result = fmt_code(code)
    return HttpResponse(status_code=200, body=json.dumps({"body": result}))


router.post("/-/play/compile", compile_handler)
router.post("/-/play/fmt", fmt_handler)

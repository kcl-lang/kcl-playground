"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.invokeKCLFmt = exports.invokeKCLRun = exports.load = void 0;
const wasi_1 = require("@wasmer/wasi");
const RUN_FUNCTION_NAME = "kcl_run";
const FMT_FUNCTION_NAME = "kcl_fmt";
const RUNTIME_ERR_FUNCTION_NAME = "kcl_runtime_err";
/**
 * load the KCL WASM
 * @param options
 * @returns
 */
async function load(opts) {
    await (0, wasi_1.init)();
    const options = opts ?? {};
    const w = new wasi_1.WASI({
        env: options.env ?? {},
        fs: options.fs,
    });
    let bytes;
    if (options.data) {
        bytes = options.data;
    }
    else {
        if (typeof window !== "undefined") {
            const response = await fetch("../kcl.wasm");
            bytes = await response.arrayBuffer();
        }
        else {
            throw new Error("Unsupported environment");
        }
    }
    const imports = {
        env: {
            kclvm_plugin_invoke_json_wasm: (_method, _args, _kwargs) => {
                // TODO: KCL WASM plugin impl
                return 0;
            },
        },
        ...(options.imports ?? {}),
    };
    const module = await WebAssembly.compile(bytes);
    // Instantiate the WASI module
    return w.instantiate(module, imports);
}
exports.load = load;
/**
 * Exported function to invoke the KCL run operation.
 */
function invokeKCLRun(instance, opts) {
    const exports = instance.exports;
    const [filenamePtr, filenamePtrLength] = copyStringToWasmMemory(instance, opts.filename);
    const [sourcePtr, sourcePtrLength] = copyStringToWasmMemory(instance, opts.source);
    let result;
    try {
        const resultPtr = exports[RUN_FUNCTION_NAME](filenamePtr, sourcePtr);
        const [resultStr, resultPtrLength] = copyCStrFromWasmMemory(instance, resultPtr);
        exports.kcl_free(resultPtr, resultPtrLength);
        result = resultStr;
    }
    catch (error) {
        let runtimeErrPtrLength = 1024;
        let runtimeErrPtr = exports.kcl_malloc(runtimeErrPtrLength);
        exports[RUNTIME_ERR_FUNCTION_NAME](runtimeErrPtr, runtimeErrPtrLength);
        const [runtimeErrStr, _] = copyCStrFromWasmMemory(instance, runtimeErrPtr);
        exports.kcl_free(runtimeErrPtr, runtimeErrPtrLength);
        result = "ERROR:" + runtimeErrStr;
    }
    finally {
        exports.kcl_free(filenamePtr, filenamePtrLength);
        exports.kcl_free(sourcePtr, sourcePtrLength);
    }
    return result;
}
exports.invokeKCLRun = invokeKCLRun;
/**
 * Exported function to invoke the KCL format operation.
 */
function invokeKCLFmt(instance, opts) {
    const exports = instance.exports;
    const [sourcePtr, sourcePtrLength] = copyStringToWasmMemory(instance, opts.source);
    const resultPtr = exports[FMT_FUNCTION_NAME](sourcePtr);
    const [resultStr, resultPtrLength] = copyCStrFromWasmMemory(instance, resultPtr);
    exports.kcl_free(sourcePtr, sourcePtrLength);
    exports.kcl_free(resultPtr, resultPtrLength);
    return resultStr;
}
exports.invokeKCLFmt = invokeKCLFmt;
function copyStringToWasmMemory(instance, str) {
    const exports = instance.exports;
    const encodedString = new TextEncoder().encode(str);
    const pointer = exports.kcl_malloc(encodedString.length + 1); // Allocate memory and get pointer
    const buffer = new Uint8Array(exports.memory.buffer, pointer, encodedString.length + 1);
    buffer.set(encodedString);
    buffer[encodedString.length] = 0; // Null-terminate the string
    return [pointer, encodedString.length + 1];
}
function copyCStrFromWasmMemory(instance, ptr) {
    const exports = instance.exports;
    const memory = new Uint8Array(exports.memory.buffer);
    let end = ptr;
    while (memory[end] !== 0) {
        end++;
    }
    const result = new TextDecoder().decode(memory.slice(ptr, end));
    return [result, end + 1 - ptr];
}
//# sourceMappingURL=index.js.map

import { createClosure } from "@plutolang/base/closure";

export default (async () => {

const kcl_playground_dev__plutolang_pluto_Website_kcl_playground = await (
  await import("@plutolang/pluto-infra")
).Website.createInstance("./website", "kcl-playground", undefined);

const kcl_playground_dev__plutolang_pluto_Router_router = await (
  await import("@plutolang/pluto-infra")
).Router.createInstance("router", undefined);

const router_0_get_1_fn_func = async (...args: any[]) => {
  const handler = (await import("/Users/lingzhi/_Code/KCLOpenSource/kcl-playground/kcl-playground/.pluto/dev/closures/router_0_get_1_fn")).default;
  return await handler(...args);
}
const router_0_get_1_fn = createClosure(router_0_get_1_fn_func, {
  dirpath: "/Users/lingzhi/_Code/KCLOpenSource/kcl-playground/kcl-playground/.pluto/dev/closures/router_0_get_1_fn",
  exportName: "_default",
  dependencies: [],
  accessedEnvVars: [],
});

kcl_playground_dev__plutolang_pluto_Router_router.get("/-/play/compile", router_0_get_1_fn);
kcl_playground_dev__plutolang_pluto_Website_kcl_playground.postProcess();
kcl_playground_dev__plutolang_pluto_Router_router.postProcess();
return {
kcl_playground_dev__plutolang_pluto_Website_kcl_playground: kcl_playground_dev__plutolang_pluto_Website_kcl_playground.outputs,
kcl_playground_dev__plutolang_pluto_Router_router: kcl_playground_dev__plutolang_pluto_Router_router.outputs
}
})();

// Moon Bridge Cloudflare Worker — TypeScript bootstrap.
import "./vendor/wasm_exec.js";

import {createRuntimeContext, loadModule, RuntimeContext} from "./runtime";

let goInstance: Go | undefined;
let wasmInstance: WebAssembly.Instance | undefined;
const binding: Record<string, unknown> = {};

async function initWasm(ctx: RuntimeContext): Promise<void> {
    const mod = await loadModule();
    goInstance = new Go();

    let ready: () => void;
    const readyPromise = new Promise<void>((resolve) => {
        ready = resolve;
    });

    (globalThis as any).context = ctx;

    wasmInstance = new WebAssembly.Instance(mod, {
        ...goInstance.importObject,
        workers: {
            ready: () => {
                ready!();
            },
        },
    });

    goInstance.run(wasmInstance, ctx);
    await readyPromise;
}

function setContext(
    env: Record<string, unknown>,
    execCtx: ExecutionContext
): void {
    (globalThis as any).context = createRuntimeContext(env, execCtx, binding);
}

async function handleRequest(req: Request): Promise<Response> {
    try {
        const handler = binding.handleRequest as ((r: Request) => unknown) | undefined;
        if (!handler) {
            return new Response("Worker not ready", {status: 503});
        }
        const result = handler(req);
        if (result && typeof (result as Promise<unknown>).then === "function") {
            return await (result as Promise<Response>);
        }
        return result as unknown as Response;
    } catch (e) {
        console.error("handleRequest error:", e);
        return new Response("internal error", {status: 500});
    }
}

// noinspection JSUnusedGlobalSymbols
export default {
    async fetch(
        req: Request,
        env: Record<string, unknown>,
        execCtx: ExecutionContext
    ): Promise<Response> {
        if (!goInstance) {
            const ctx = createRuntimeContext(env, execCtx, binding);
            await initWasm(ctx);
        } else {
            setContext(env, execCtx);
        }
        return handleRequest(req);
    },

};

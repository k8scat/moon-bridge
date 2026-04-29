// Runtime helpers for Cloudflare Workers Go WASM integration.
// Compiled to build/runtime.mjs — paths are relative to the build/ directory.

import { connect } from "cloudflare:sockets";
import mod from "../../build/app.wasm";

export async function loadModule(): Promise<WebAssembly.Module> {
  return mod as unknown as WebAssembly.Module;
}

export interface RuntimeContext {
  env: Record<string, unknown>;
  ctx: ExecutionContext;
  connect: typeof connect;
  binding: Record<string, unknown>;
}

export function createRuntimeContext(
  env: Record<string, unknown>,
  ctx: ExecutionContext,
  binding: Record<string, unknown>
): RuntimeContext {
  return { env, ctx, connect, binding };
}

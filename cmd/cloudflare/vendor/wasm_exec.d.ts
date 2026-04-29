// Type declarations for TinyGo's wasm_exec.js (vendored, side-effect only).
// The file defines the Go class on globalThis at module evaluation time.
declare global {
    declare class Go {
        importObject: WebAssembly.Imports;
        exited: boolean;
        exitCode: number;
        _resolveExitPromise: () => void;
        _inst: WebAssembly.Instance;
        _values: unknown[];
        _goRefCounts: number[];
        _ids: Map<unknown, number>;
        _idPool: number[];
        _scheduledTimeouts: Map<number, unknown>;
        _nextCallbackTimeoutID: number;
        mem: DataView;

        constructor();

        run(instance: WebAssembly.Instance, context?: unknown): Promise<void>;

        _resume(): Promise<void>;

        _makeFuncWrapper(id: number): (...args: unknown[]) => void;
    }
}

export {};

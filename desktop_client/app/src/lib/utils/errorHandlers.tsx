import { getLogger } from "../grpcClient/client";

export async function must<T>(x: Promise<T>) {
  try {
    return await x;
  } catch (error: unknown) {
    getLogger().error("Must call failed", {
      error: JSON.stringify(error),
      mustCaller: must.caller.name,
    });

    throw error;
  }
}

export function mustVoid<T>(x: void | Promise<T>) {
  if (!x) {
    return;
  }

  must(x);
}

import { FinishedUnaryCall } from "@protobuf-ts/runtime-rpc";

export function newResult<T extends Object, V extends Object>(
  req: T,
  resp: V,
): FinishedUnaryCall<T, V> {
  return {
    response: resp,
    request: req,
    method: {
      name: "Mocked call",
      clientStreaming: false,
      I: {} as any,
      O: {} as any,
      idempotency: "IDEMPOTENT",
      localName: "Mocked call",
      options: {},
      serverStreaming: false,
      service: {
        options: {},
        methods: [],
        typeName: "Mocked call",
      },
    },
    headers: {},
    status: 200 as any,
    requestHeaders: {},
    trailers: {},
  };
}

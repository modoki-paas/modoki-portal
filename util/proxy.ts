interface ProxiedRequest {
  url: string;
  method: string;
  headers?: {[key: string]: string};
  body?: string;
}

export async function fetch(input: RequestInfo, init?: RequestInit): Promise<Response> {
  if (typeof input === "string") {
    const reqBody = {
      url: input,
      method: init?.method ?? "GET",
      headers: init?.headers,
      body: init?.body?.toString(),
    } as ProxiedRequest;

    return await window.fetch(
      "http://localhost:8080/proxy",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(reqBody),
        credentials: "include"
      },
    )
  }

  throw "Request type for input is not supported"
}

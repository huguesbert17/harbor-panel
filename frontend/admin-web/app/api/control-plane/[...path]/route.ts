import { NextResponse } from "next/server";

const defaultApiBaseUrl = "http://localhost:8080";

function apiBaseUrl() {
  return (process.env.HARBOR_API_BASE_URL ?? defaultApiBaseUrl).replace(/\/$/, "");
}

export async function GET(request: Request, context: { params: Promise<{ path: string[] }> }) {
  const { path } = await context.params;
  const endpoint = path.join("/");
  if (!isAllowedEndpoint(endpoint)) {
    return NextResponse.json({ error: "Not found" }, { status: 404 });
  }

  const upstreamUrl = new URL(`${apiBaseUrl()}/api/v1/${endpoint}`);
  upstreamUrl.search = new URL(request.url).search;

  try {
    const response = await fetch(upstreamUrl, { headers: { Accept: "application/json" }, cache: "no-store" });
    const body = await response.text();
    return new NextResponse(body, { status: response.status, headers: { "content-type": response.headers.get("content-type") ?? "application/json" } });
  } catch {
    return NextResponse.json({ error: "The control-plane API could not be reached." }, { status: 502 });
  }
}

export async function POST(request: Request, context: { params: Promise<{ path: string[] }> }) {
  const { path } = await context.params;
  const endpoint = path.join("/");
  if (!isAllowedEndpoint(endpoint)) return NextResponse.json({ error: "Not found" }, { status: 404 });

  const upstreamUrl = new URL(`${apiBaseUrl()}/api/v1/${endpoint}`);
  upstreamUrl.search = new URL(request.url).search;
  try {
    const response = await fetch(upstreamUrl, {
      method: "POST",
      headers: { Accept: "application/json", "content-type": request.headers.get("content-type") ?? "application/json" },
      body: await request.text(),
      cache: "no-store",
    });
    const body = await response.text();
    return new NextResponse(body, { status: response.status, headers: { "content-type": response.headers.get("content-type") ?? "application/json" } });
  } catch {
    return NextResponse.json({ error: "The control-plane API could not be reached." }, { status: 502 });
  }
}

function isAllowedEndpoint(endpoint: string) {
  return endpoint === "accounts" || endpoint === "servers" || endpoint === "domains" || endpoint === "dns/zones" || /^dns\/zones\/[^/]+\/records$/.test(endpoint);
}


import { NextResponse } from 'next/server'

// Proxy para evitar CORS / mixed-content: el browser llama a same-origin (/api/demo/login)
// y el servidor Next.js llama al backend.
function backendBase(): string {
  return (
    process.env.BACKEND_URL ||
    process.env.NEXT_PUBLIC_BACKEND_URL ||
    'http://localhost:8080'
  ).replace(/\/$/, '')
}

export async function POST(req: Request) {
  const body = await req.json().catch(() => ({}))
  const url = `${backendBase()}/api/v1/demo/login`

  try {
    const resp = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
      cache: 'no-store',
    })

    const text = await resp.text()
    return new NextResponse(text, {
      status: resp.status,
      headers: {
        'Content-Type': resp.headers.get('content-type') || 'application/json',
      },
    })
  } catch (e: any) {
    return NextResponse.json(
      {
        error: 'No se pudo conectar al backend',
        detail: String(e?.message || e),
        backend_url: url,
      },
      { status: 502 }
    )
  }
}

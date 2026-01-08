import { NextResponse } from 'next/server'

function backendBase(): string {
  return (
    process.env.BACKEND_URL ||
    process.env.NEXT_PUBLIC_BACKEND_URL ||
    'http://localhost:8080'
  ).replace(/\/$/, '')
}

async function proxy(req: Request, ctx: { params: Promise<{ path?: string[] }> }) {
  const { path } = await ctx.params
  const tail = (path || []).join('/')
  const url = `${backendBase()}/api/v1/${tail}`

  let body: string | undefined
  if (req.method !== 'GET' && req.method !== 'HEAD') {
    body = await req.text()
  }

  try {
    const resp = await fetch(url, {
      method: req.method,
      headers: {
        'Content-Type': req.headers.get('content-type') || 'application/json',
      },
      body,
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

export async function GET(req: Request, ctx: any) {
  return proxy(req, ctx)
}

export async function POST(req: Request, ctx: any) {
  return proxy(req, ctx)
}

export async function PUT(req: Request, ctx: any) {
  return proxy(req, ctx)
}

export async function DELETE(req: Request, ctx: any) {
  return proxy(req, ctx)
}

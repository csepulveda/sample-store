import "@/otel"
import type { NextApiRequest, NextApiResponse } from "next"
import { context, trace, propagation } from "@opentelemetry/api"

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const tracer = trace.getTracer("ui-web")
  const baseUrl = process.env.PRODUCT_API_BASE_URL || "http://products-service:8080"

  return tracer.startActiveSpan(
    "proxy_product_request",
    {},
    context.active(),
    async (span) => {
      try {
        if (req.method !== "GET") {
          span.setStatus({ code: 1, message: "Method not allowed" })
          return res.status(405).end()
        }

        const headers: Record<string, string> = {}
        propagation.inject(context.active(), headers)

        const response = await fetch(`${baseUrl}/api/products`, {
          method: "GET",
          headers,
        })

        const data = await response.json()
        span.setStatus({ code: 0 })
        return res.status(response.status).json(data)
      } catch (error) {
        span.recordException(error as Error)
        span.setStatus({ code: 2, message: "Unexpected error" })
        return res.status(500).json({ error: "Internal Server Error" })
      } finally {
        span.end()
      }
    }
  )
}
import "@/otel"
import type { NextApiRequest, NextApiResponse } from "next"
import { context, trace, propagation } from "@opentelemetry/api"

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const tracer = trace.getTracer("ui-web")
  const baseUrl = process.env.ORDER_API_BASE_URL || "http://orders-service:8080"

  return tracer.startActiveSpan(
    "proxy_order_request",
    {},
    context.active(),
    async (span) => {
      try {
        if (req.method !== "POST") {
          span.setStatus({ code: 1, message: "Method not allowed" })
          return res.status(405).end()
        }

        const { items } = req.body

        if (!items || !Array.isArray(items) || items.length === 0) {
          span.setStatus({ code: 1, message: "Missing or invalid items" })
          return res.status(400).json({ error: "Order must have at least one item" })
        }

        const headers: Record<string, string> = {
          "Content-Type": "application/json",
        }

        propagation.inject(context.active(), headers)

        const response = await fetch(`${baseUrl}/api/orders`, {
          method: "POST",
          headers,
          body: JSON.stringify({ items }),
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
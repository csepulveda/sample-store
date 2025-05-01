import "@/otel"
import type { NextApiRequest, NextApiResponse } from "next"
import { context, trace, propagation } from "@opentelemetry/api"

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const tracer = trace.getTracer("ui-backoffice")
  const { id } = req.query
  const baseUrl = process.env.ORDER_API_BASE_URL || "http://localhost:8080"

  return tracer.startActiveSpan("proxy_order_id", {}, context.active(), async (span) => {
    try {
      const headers: Record<string, string> = { "Content-Type": "application/json" }
      propagation.inject(context.active(), headers)

      if (req.method === "PATCH") {
        const response = await fetch(`${baseUrl}/api/orders/${id}`, {
          method: "PATCH",
          headers,
          body: JSON.stringify(req.body),
        })
        const data = await response.json()
        span.setStatus({ code: 0 })
        return res.status(response.status).json(data)
      }

      if (req.method === "DELETE") {
        const response = await fetch(`${baseUrl}/api/orders/${id}`, {
          method: "DELETE",
          headers,
        })
        span.setStatus({ code: 0 })
        return res.status(response.status).end()
      }

      span.setStatus({ code: 1, message: "Method not allowed" })
      return res.status(405).end()
    } catch (error) {
      span.recordException(error as Error)
      span.setStatus({ code: 2, message: "Unexpected error" })
      return res.status(500).json({ error: "Internal Server Error" })
    } finally {
      span.end()
    }
  })
}
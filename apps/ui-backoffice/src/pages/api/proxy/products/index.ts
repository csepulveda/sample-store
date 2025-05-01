import "@/otel"
import type { NextApiRequest, NextApiResponse } from "next"
import { context, trace, propagation } from "@opentelemetry/api"

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const tracer = trace.getTracer("ui-backoffice")
  const baseUrl = process.env.PRODUCT_API_BASE_URL || "http://products-service:8080"

  return tracer.startActiveSpan("proxy_product_list", {}, context.active(), async (span) => {
    try {
      const headers: Record<string, string> = { "Content-Type": "application/json" }
      propagation.inject(context.active(), headers)

      if (req.method === "GET") {
        const response = await fetch(`${baseUrl}/api/products`, { headers })
        const data = await response.json()
        span.setStatus({ code: 0 })
        return res.status(response.status).json(data)
      }

      if (req.method === "POST") {
        const { name, price, stock } = req.body

        if (!name || typeof name !== "string" || name.trim() === "") {
          span.setStatus({ code: 1, message: "Invalid name" })
          return res.status(400).json({ error: "Name is required" })
        }

        if (price == null || typeof price !== "number" || price <= 0) {
          span.setStatus({ code: 1, message: "Invalid price" })
          return res.status(400).json({ error: "Price must be greater than 0" })
        }

        if (stock == null || typeof stock !== "number" || stock < 0) {
          span.setStatus({ code: 1, message: "Invalid stock" })
          return res.status(400).json({ error: "Stock must be zero or greater" })
        }

        const response = await fetch(`${baseUrl}/api/products`, {
          method: "POST",
          headers,
          body: JSON.stringify(req.body),
        })
        const data = await response.json()
        span.setStatus({ code: 0 })
        return res.status(response.status).json(data)
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
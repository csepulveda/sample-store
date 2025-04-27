import type { NextApiRequest, NextApiResponse } from "next"

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const baseUrl = process.env.ORDER_API_BASE_URL || "http://orders-service:8080"

  if (req.method === "POST") {
    const { items } = req.body

    if (!items || !Array.isArray(items) || items.length === 0) {
      return res.status(400).json({ error: "Order must have at least one item" })
    }

    const response = await fetch(`${baseUrl}/api/orders`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ items }),
    })

    const data = await response.json()
    return res.status(response.status).json(data)
  }

  return res.status(405).end()
}
import type { NextApiRequest, NextApiResponse } from "next"

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const baseUrl = process.env.PRODUCT_API_BASE_URL || "http://products-service:8080"

  if (req.method === "GET") {
    const response = await fetch(`${baseUrl}/api/products`)
    const data = await response.json()
    return res.status(response.status).json(data)
  }

  if (req.method === "POST") {
    const { name, price, stock } = req.body

    if (!name || typeof name !== "string" || name.trim() === "") {
      return res.status(400).json({ error: "Name is required" })
    }

    if (price == null || typeof price !== "number" || price <= 0) {
      return res.status(400).json({ error: "Price must be greater than 0" })
    }

    if (stock == null || typeof stock !== "number" || stock < 0) {
      return res.status(400).json({ error: "Stock must be zero or greater" })
    }

    const response = await fetch(`${baseUrl}/api/products`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req.body),
    })
    const data = await response.json()
    return res.status(response.status).json(data)
  }

  return res.status(405).end()
}
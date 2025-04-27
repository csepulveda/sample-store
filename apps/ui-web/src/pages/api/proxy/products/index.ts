import type { NextApiRequest, NextApiResponse } from "next"

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const baseUrl = process.env.PRODUCT_API_BASE_URL || "http://products-service:8080"

  if (req.method === "GET") {
    const response = await fetch(`${baseUrl}/api/products`)
    const data = await response.json()
    return res.status(response.status).json(data)
  }

  return res.status(405).end()
}
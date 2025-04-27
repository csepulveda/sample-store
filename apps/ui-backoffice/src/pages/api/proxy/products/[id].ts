import type { NextApiRequest, NextApiResponse } from "next"

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const { id } = req.query
  const baseUrl = process.env.PRODUCT_API_BASE_URL || "http://localhost:8080"
  console.log("product_baseurl", baseUrl)

  if (req.method === "PATCH") {
    const response = await fetch(`${baseUrl}/api/products/${id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req.body),
    })
    const data = await response.json()
    return res.status(response.status).json(data)
  }

  if (req.method === "DELETE") {
    const response = await fetch(`${baseUrl}/api/products/${id}`, {
      method: "DELETE",
    })
    return res.status(response.status).end()
  }

  return res.status(405).end()
}
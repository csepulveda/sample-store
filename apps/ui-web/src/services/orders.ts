import { Product } from "./products"

export interface OrderItem {
  productId: string
  quantity: number
}

export interface CreateOrderPayload {
  items: OrderItem[]
}

export async function createOrder(payload: CreateOrderPayload) {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/orders`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })

  if (!res.ok) {
    throw new Error("Failed to create order")
  }

  return res.json()
}
export interface OrderItem {
  productId: string
  quantity: number
}

export async function createOrder(items: OrderItem[]) {
  const res = await fetch("/api/proxy/orders", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ items }),
  })
  if (!res.ok) {
    const error = await res.json()
    throw new Error(error.error || "Failed to create order")
  }
}
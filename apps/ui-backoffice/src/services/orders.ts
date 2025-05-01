export interface OrderItem {
  productId: string
  productName: string
  quantity: number
}

export interface Order {
  id: string
  status: string
  createdAt: string
  items: OrderItem[]
}

function getBaseUrl() {
  if (typeof window !== "undefined") {
      return "" // Client-side uses relative path
  }
  return process.env.NEXT_PUBLIC_APP_URL || "http://ui-backoffice:3000"
}

export async function fetchOrdersServerSide(): Promise<Order[]> {
  const baseUrl = getBaseUrl()
  const res = await fetch(`${baseUrl}/api/proxy/orders`)
  if (!res.ok) {
    throw new Error("Failed to fetch orders")
  }
  const data = await res.json()
  return Array.isArray(data) ? data : []
}

export async function updateOrderStatus(id: string, status: string) {
  const baseUrl = getBaseUrl()
  const res = await fetch(`${baseUrl}/api/proxy/orders/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ status }),
  })
  if (!res.ok) {
    throw new Error("Failed to update order status")
  }
}

export async function deleteOrder(id: string) {
  const baseUrl = getBaseUrl()
  const res = await fetch(`${baseUrl}/api/proxy/orders/${id}`, {
    method: "DELETE",
  })
  if (!res.ok) {
    throw new Error("Failed to delete order")
  }
}
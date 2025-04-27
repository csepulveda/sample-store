export interface OrderItem {
    productId: string
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
      return "" // Browser uses relative URL
    }
    return process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000"
  }
  
  export async function fetchOrdersServerSide(): Promise<Order[]> {
    const baseUrl = getBaseUrl()
    const res = await fetch(`${baseUrl}/api/proxy/orders`)
    if (!res.ok) {
      throw new Error("Failed to fetch orders")
    }
    return res.json()
  }
  
  export async function createOrder(payload: { items: OrderItem[] }) {
    const baseUrl = getBaseUrl()
    const res = await fetch(`${baseUrl}/api/proxy/orders`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    })
    if (!res.ok) {
      throw new Error("Failed to create order")
    }
  }
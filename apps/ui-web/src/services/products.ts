export interface Product {
  id: string
  name: string
  description: string
  price: number
  stock: number
}

function getBaseUrl() {
  if (typeof window !== "undefined") {
    return "" // Cliente puede usar relative path
  }
  return process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000"
}

export async function fetchProductsServerSide(): Promise<Product[]> {
  const baseUrl = getBaseUrl()
  const res = await fetch(`${baseUrl}/api/proxy/products`)
  if (!res.ok) {
    throw new Error("Failed to fetch products")
  }
  return res.json()
}
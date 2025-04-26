export interface Product {
    id: string
    name: string
    description: string
    price: number
    stock: number
  }
  
  export async function fetchProducts(): Promise<Product[]> {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/products`)
    if (!res.ok) {
      throw new Error('Failed to fetch products')
    }
    return res.json()
  }
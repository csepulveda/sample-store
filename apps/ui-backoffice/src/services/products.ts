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
  
  export async function deleteProduct(id: string) {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/products/${id}`, {
      method: "DELETE",
    })
    if (!res.ok) {
      throw new Error('Failed to delete product')
    }
  }
  
  export async function updateProduct(id: string, product: Partial<Product>) {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/products/${id}`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(product),
    })
    if (!res.ok) {
      throw new Error('Failed to update product')
    }
  }
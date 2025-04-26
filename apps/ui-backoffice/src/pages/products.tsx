import { useEffect, useState } from "react"
import { fetchProducts, Product } from "@/services/products"
import { ProductTable } from "@/components/ProductTable"

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([])

  const loadProducts = async () => {
    const data = await fetchProducts()
    setProducts(data)
  }

  useEffect(() => {
    loadProducts()
  }, [])

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8 text-white">Manage Products</h1>
      <ProductTable products={products} reload={loadProducts} />
    </div>
  )
}
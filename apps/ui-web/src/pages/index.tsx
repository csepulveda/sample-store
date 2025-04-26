import { useEffect, useState } from "react"
import { fetchProducts, Product } from "@/services/products"
import { ProductCard } from "@/components/ProductCard"
import { CartSummary } from "@/components/CartSummary"
import { CartProvider } from "@/context/cart"

export default function Home() {
  const [products, setProducts] = useState<Product[]>([])

  useEffect(() => {
    fetchProducts().then(setProducts)
  }, [])

  return (
    <CartProvider>
      <div className="container mx-auto p-8">
        <h1 className="text-3xl font-bold mb-8 text-white">Products</h1>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {products.map(product => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
        <CartSummary />
      </div>
    </CartProvider>
  )
}
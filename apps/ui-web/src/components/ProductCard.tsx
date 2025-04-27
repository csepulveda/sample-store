"use client"

import { Product } from "@/services/products"
import { useCart } from "@/context"

export function ProductCard({ product }: { product: Product }) {
  const { dispatch } = useCart()

  const handleAddToCart = () => {
    console.log("CLICK on product:", product.name)
    dispatch({ type: "ADD_TO_CART", payload: { productId: product.id, quantity: 1 } })
  }

  return (
    <div className="bg-zinc-800 p-4 rounded-lg flex flex-col justify-between">
      <div>
        <h3 className="text-xl font-bold text-white mb-2">{product.name}</h3>
        <p className="text-white text-sm mb-4">{product.description}</p>
        <p className="text-white font-bold">${product.price.toFixed(2)}</p>
      </div>
      <button
        onClick={handleAddToCart}
        className="mt-4 bg-green-600 hover:bg-green-500 text-white py-2 rounded"
      >
        Add to Cart
      </button>
    </div>
  )
}
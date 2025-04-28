"use client"

import { Product } from "@/services/products"
import { useCart } from "@/context"

export function ProductCard({ product }: { product: Product }) {
  const { dispatch, cartItems } = useCart()

  const quantityInCart = cartItems.find(item => item.productId === product.id)?.quantity || 0
  const availableStock = product.stock - quantityInCart

  const handleAddToCart = () => {
    if (availableStock <= 0) return
    dispatch({ type: "ADD_TO_CART", payload: { productId: product.id, quantity: 1 } })
  }

  return (
    <div className="bg-zinc-800 p-4 rounded-lg flex flex-col justify-between">
      <div>
        <h3 className="text-xl font-bold text-white mb-2">{product.name}</h3>
        <p className="text-white text-sm mb-2">{product.description}</p>
        <p className="text-white font-bold mb-2">${product.price.toFixed(2)}</p>
        <p className="text-white text-sm mb-2">Stock: {availableStock}</p>
      </div>
      <button
        onClick={handleAddToCart}
        disabled={availableStock <= 0}
        className={`mt-4 ${
          availableStock <= 0 ? "bg-gray-500 cursor-not-allowed" : "bg-green-600 hover:bg-green-500"
        } text-white py-2 rounded`}
      >
        {availableStock <= 0 ? "Out of Stock" : "Add to Cart"}
      </button>
    </div>
  )
}
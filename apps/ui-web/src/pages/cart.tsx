"use client"

import { useCart } from "@/context"

export default function CartPage() {
  const { cartItems } = useCart()

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8 text-white">Your Cart</h1>

      {cartItems.length === 0 ? (
        <p className="text-white">Your cart is empty.</p>
      ) : (
        <div className="flex flex-col gap-4">
          {cartItems.map((item, idx) => (
            <div key={idx} className="bg-zinc-800 p-4 rounded">
              <p className="text-white">Product ID: {item.productId}</p>
              <p className="text-white">Quantity: {item.quantity}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
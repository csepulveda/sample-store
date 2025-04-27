"use client"

import { useCart } from "@/context"
import { createOrder } from "@/services/orders"
import { useState } from "react"

export function CartSummary() {
  const { cartItems, dispatch } = useCart()
  const [loading, setLoading] = useState(false)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

  const handleCheckout = async () => {
    if (cartItems.length === 0) {
      setErrorMessage("Your cart is empty.")
      return
    }

    setLoading(true)
    setErrorMessage(null)
    setSuccessMessage(null)

    try {
      await createOrder(cartItems)
      dispatch({ type: "CLEAR_CART" })
      setSuccessMessage("Order created successfully!")
    } catch (err: any) {
      setErrorMessage(err.message || "Failed to create order")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="bg-zinc-800 p-6 rounded-lg">
      <h2 className="text-2xl font-bold mb-4 text-white">Cart Summary</h2>

      {cartItems.length === 0 ? (
        <p className="text-white mb-4">Your cart is empty.</p>
      ) : (
        <ul className="mb-4">
          {cartItems.map((item, idx) => (
            <li key={idx} className="text-white">
              {item.productId} - Qty: {item.quantity}
            </li>
          ))}
        </ul>
      )}

      {errorMessage && <p className="text-red-500 mb-4">{errorMessage}</p>}
      {successMessage && <p className="text-green-500 mb-4">{successMessage}</p>}

      <button
        onClick={handleCheckout}
        disabled={loading || cartItems.length === 0}
        className={`w-full ${
          loading ? "bg-gray-600" : "bg-blue-600 hover:bg-blue-500"
        } text-white font-semibold py-2 rounded`}
      >
        {loading ? "Processing..." : "Checkout"}
      </button>
    </div>
  )
}
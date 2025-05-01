"use client"

import { useCart } from "@/context"
import { createOrder } from "@/services/orders"
import { useState } from "react"

export default function CartPage() {
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
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8 text-white">Your Cart</h1>

      {cartItems.length === 0 ? (
        <p className="text-white">Your cart is empty.</p>
      ) : (
        <div className="flex flex-col gap-4 mb-8">
          {cartItems.map((item, idx) => (
            <div key={idx} className="bg-zinc-800 p-4 rounded">
              <p className="text-white">Product Name: {item.productName}</p>
              <p className="text-white">Quantity: {item.quantity}</p>
            </div>
          ))}
        </div>
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
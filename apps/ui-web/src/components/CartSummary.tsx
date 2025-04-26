import { useCart } from "@/context/cart"
import { createOrder } from "@/services/orders"

export function CartSummary() {
  const { state, dispatch } = useCart()

  const total = state.items.reduce(
    (sum, item) => sum + item.price * item.quantity,
    0
  )

  const handleCheckout = async () => {
    try {
      const payload = {
        items: state.items.map(item => ({
          productId: item.id,
          quantity: item.quantity,
        })),
      }
      await createOrder(payload)
      dispatch({ type: "CLEAR_CART" })
      alert("Order created successfully")
    } catch (err) {
      console.error(err)
      alert("Error creating order")
    }
  }

  return (
    <div className="bg-zinc-800 text-white p-4 rounded-lg shadow-md mt-8">
      <h2 className="text-lg font-semibold mb-4">Cart</h2>
      {state.items.length === 0 ? (
        <p className="text-gray-400">Cart is empty</p>
      ) : (
        <>
          <ul className="space-y-2">
            {state.items.map(item => (
              <li key={item.id} className="flex justify-between">
                <span>
                  {item.name} x {item.quantity}
                </span>
                <span>${(item.price * item.quantity).toFixed(2)}</span>
              </li>
            ))}
          </ul>
          <div className="mt-4 flex justify-between font-bold border-t border-zinc-600 pt-2">
            <span>Total</span>
            <span>${total.toFixed(2)}</span>
          </div>
          <div className="flex gap-2 mt-4">
            <button
              onClick={() => dispatch({ type: "CLEAR_CART" })}
              className="bg-red-600 hover:bg-red-500 text-white px-4 py-2 rounded"
            >
              Clear cart
            </button>
            <button
              onClick={handleCheckout}
              className="bg-green-600 hover:bg-green-500 text-white px-4 py-2 rounded"
            >
              Checkout
            </button>
          </div>
        </>
      )}
    </div>
  )
}
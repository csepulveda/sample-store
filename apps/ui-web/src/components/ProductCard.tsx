import { Product } from "@/services/products"
import { useCart } from "@/context/cart"

export function ProductCard({ product }: { product: Product }) {
  const { dispatch } = useCart()

  return (
    <div className="border rounded-lg p-4 shadow-md flex flex-col bg-zinc-900">
      <h2 className="text-xl font-semibold text-white">{product.name}</h2>
      <p className="text-gray-400">{product.description}</p>
      <p className="mt-2 font-bold text-white">${product.price}</p>
      <p className="text-sm text-gray-500">Stock: {product.stock}</p>
      <button
        onClick={() => dispatch({ type: "ADD_ITEM", product })}
        className="mt-4 bg-indigo-600 text-white py-2 px-4 rounded hover:bg-indigo-500"
      >
        Add to cart
      </button>
    </div>
  )
}
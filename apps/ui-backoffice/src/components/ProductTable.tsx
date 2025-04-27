"use client"

import { Product, deleteProduct, updateProduct } from "@/services/products"
import { useState } from "react"

export function ProductTable({ products, reload }: { products: Product[]; reload: () => void }) {
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editData, setEditData] = useState<Partial<Product>>({})

  const handleEdit = (product: Product) => {
    setEditingId(product.id)
    setEditData({
      name: product.name,
      description: product.description,
      price: product.price,
      stock: product.stock,
    })
  }

  const handleSave = async (id: string) => {
    await updateProduct(id, editData)
    setEditingId(null)
    reload()
  }

  const handleDelete = async (id: string) => {
    await deleteProduct(id)
    reload()
  }

  const handleChange = (field: keyof Product, value: string | number) => {
    setEditData(prev => ({ ...prev, [field]: value }))
  }

  return (
    <table className="min-w-full bg-zinc-900 text-white rounded-lg overflow-hidden">
      <thead className="bg-zinc-700">
        <tr>
          <th className="py-2 px-4 text-left">Name</th>
          <th className="py-2 px-4 text-left">Description</th>
          <th className="py-2 px-4 text-left">Price</th>
          <th className="py-2 px-4 text-left">Stock</th>
          <th className="py-2 px-4 text-left">Actions</th>
        </tr>
      </thead>
      <tbody>
        {products.map(product => (
          <tr key={product.id} className="border-t border-zinc-700">
            <td className="py-2 px-4">
              {editingId === product.id ? (
                <input
                  type="text"
                  value={editData.name || ""}
                  onChange={(e) => handleChange("name", e.target.value)}
                  className="bg-zinc-800 border border-zinc-600 rounded px-2 py-1 w-full"
                />
              ) : (
                product.name
              )}
            </td>
            <td className="py-2 px-4">
              {editingId === product.id ? (
                <input
                  type="text"
                  value={editData.description || ""}
                  onChange={(e) => handleChange("description", e.target.value)}
                  className="bg-zinc-800 border border-zinc-600 rounded px-2 py-1 w-full"
                />
              ) : (
                product.description
              )}
            </td>
            <td className="py-2 px-4">
              {editingId === product.id ? (
                <input
                  type="number"
                  value={editData.price ?? 0}
                  onChange={(e) => handleChange("price", parseFloat(e.target.value))}
                  className="bg-zinc-800 border border-zinc-600 rounded px-2 py-1 w-24"
                />
              ) : (
                `$${product.price}`
              )}
            </td>
            <td className="py-2 px-4">
              {editingId === product.id ? (
                <input
                  type="number"
                  value={editData.stock ?? 0}
                  onChange={(e) => handleChange("stock", parseInt(e.target.value))}
                  className="bg-zinc-800 border border-zinc-600 rounded px-2 py-1 w-20"
                />
              ) : (
                product.stock
              )}
            </td>
            <td className="py-2 px-4 flex gap-2">
              {editingId === product.id ? (
                <button
                  onClick={() => handleSave(product.id)}
                  className="bg-green-600 hover:bg-green-500 px-3 py-1 rounded"
                >
                  Save
                </button>
              ) : (
                <button
                  onClick={() => handleEdit(product)}
                  className="bg-blue-600 hover:bg-blue-500 px-3 py-1 rounded"
                >
                  Edit
                </button>
              )}
              <button
                onClick={() => handleDelete(product.id)}
                className="bg-red-600 hover:bg-red-500 px-3 py-1 rounded"
              >
                Delete
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
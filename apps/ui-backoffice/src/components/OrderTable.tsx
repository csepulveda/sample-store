"use client"

import { useState } from "react"
import { Order, updateOrderStatus, deleteOrder } from "@/services/orders"
import { toast } from "react-hot-toast"

export function OrderTable({ orders, reload }: { orders: Order[]; reload: () => void }) {
  const [filterStatus, setFilterStatus] = useState<string>("all")
  const [saving, setSaving] = useState(false)

  const handleStatusChange = async (id: string, newStatus: string) => {
    setSaving(true)
    try {
      await updateOrderStatus(id, newStatus)
      toast.success("Order status updated")
      reload()
    } catch (err) {
      console.error("Failed to update order status:", err)
      toast.error("Failed to update order")
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    setSaving(true)
    try {
      await deleteOrder(id)
      toast.success("Order deleted")
      reload()
    } catch (err) {
      console.error("Failed to delete order:", err)
      toast.error("Failed to delete order")
    } finally {
      setSaving(false)
    }
  }

  const filteredOrders = filterStatus === "all"
    ? orders
    : orders.filter(order => order.status === filterStatus)

  const getStatusColorClass = (status: string) => {
    switch (status) {
      case "created":
        return "bg-blue-600"
      case "shipped":
        return "bg-yellow-400 text-black"
      case "delivered":
        return "bg-green-600"
      case "returned":
        return "bg-purple-600"
      case "canceled":
        return "bg-red-600"
      default:
        return "bg-zinc-800"
    }
  }

  return (
    <div className="overflow-x-auto bg-zinc-950 p-6 rounded-lg">
      <div className="mb-6 flex gap-4 items-center">
        <label className="text-white font-semibold">Filter by Status:</label>
        <select
          value={filterStatus}
          onChange={e => setFilterStatus(e.target.value)}
          className="bg-zinc-700 text-white rounded p-2"
        >
          <option value="all">All</option>
          <option value="created">Created</option>
          <option value="shipped">Shipped</option>
          <option value="delivered">Delivered</option>
          <option value="returned">Returned</option>
          <option value="canceled">Canceled</option>
        </select>
      </div>

      <table className="min-w-full text-white rounded-lg overflow-hidden">
        <thead className="bg-gradient-to-r from-zinc-800 to-zinc-700">
          <tr>
            <th className="py-3 px-4 text-left">Order ID</th>
            <th className="py-3 px-4 text-left">Status</th>
            <th className="py-3 px-4 text-left">Created At</th>
            <th className="py-3 px-4 text-left">Items</th>
            <th className="py-3 px-4 text-left">Actions</th>
          </tr>
        </thead>
        <tbody>
          {filteredOrders.map(order => (
            <tr key={order.id} className="border-t border-zinc-600 hover:bg-zinc-800">
              <td className="py-2 px-4">{order.id}</td>
              <td className="py-2 px-4">
                <select
                  value={order.status}
                  onChange={e => handleStatusChange(order.id, e.target.value)}
                  className={`border border-zinc-600 rounded px-2 py-1 w-full text-white ${getStatusColorClass(order.status)}`}
                  disabled={saving || order.status === "canceled"} // Optional: disable canceled state
                >
                  <option value="created">Created</option>
                  <option value="shipped">Shipped</option>
                  <option value="delivered">Delivered</option>
                  <option value="returned">Returned</option>
                  <option value="canceled">Canceled</option>
                </select>
              </td>
              <td className="py-2 px-4">{new Date(order.createdAt).toLocaleString()}</td>
              <td className="py-2 px-4">
                <ul className="list-disc list-inside">
                  {order.items.map((item, idx) => (
                    <li key={idx}>
                      {item.productName} - Qty: {item.quantity}
                    </li>
                  ))}
                </ul>
              </td>
              <td className="py-2 px-4 flex gap-2">
                <button
                  onClick={() => handleDelete(order.id)}
                  disabled={saving}
                  className="bg-red-600 hover:bg-red-500 transition px-3 py-1 rounded"
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
"use client"

import { useEffect, useState } from "react"
import { Order, fetchOrdersServerSide } from "@/services/orders"
import { OrderTable } from "@/components/OrderTable"

export function OrderManager({ initialOrders }: { initialOrders: Order[] }) {
  const [orders, setOrders] = useState<Order[] | null>(null)

  useEffect(() => {
    setOrders(initialOrders)
  }, [initialOrders])

  const reload = async () => {
    try {
      const updatedOrders = await fetchOrdersServerSide()
      setOrders(updatedOrders)
    } catch (error) {
      console.error("Failed to reload orders:", error)
    }
  }

  if (orders === null) {
    return <div className="text-white">Loading...</div>
  }

  return (
    <OrderTable orders={orders} reload={reload} />
  )
}
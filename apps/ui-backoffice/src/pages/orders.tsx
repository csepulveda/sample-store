import { fetchOrdersServerSide, Order } from "@/services/orders"
import { OrderManager } from "@/components/OrderManager"

export default function OrdersPage({ initialOrders }: { initialOrders: Order[] }) {
  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8">Manage Orders</h1>
      <OrderManager initialOrders={initialOrders} />
    </div>
  )
}

export async function getServerSideProps() {
  try {
    const initialOrders = await fetchOrdersServerSide()
    return { props: { initialOrders } }
  } catch (error) {
    console.error("Failed to fetch orders:", error)
    return { props: { initialOrders: [] } }
  }
}
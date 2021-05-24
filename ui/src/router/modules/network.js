import Layout from "@/business/app-layout/horizontal-layout"

const Network = {
  path: "/network",
  sort: 3,
  component: Layout,
  name: "Network",
  meta: {
    title: "business.network.network",
    icon: "iconfont iconnetwork"
  },
  children: [
    {
      path: "/services",
      component: () => import("@/business/network/services"),
      name: "Services",
      meta: {
        title: "Services",
      }
    },
    {
      path: "/ingresses",
      component: () => import("@/business/network/ingresses"),
      name: "Ingresses",
      meta: {
        title: "Ingresses",
      }
    }
  ]
}

export default Network

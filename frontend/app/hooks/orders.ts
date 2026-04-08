import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchWithAuth } from "../utils/api";
import { type Order } from "../schemas/order";
import { type Page } from "../schemas/api";

export function useGetPage(offset: number = 0, limit: number = 10) {
  return useQuery<Page<Order>>({
    queryKey: ["orders", "page", { offset, limit }],
    queryFn: () => fetchWithAuth(`/orders?offset=${offset}&limit=${limit}`),
  });
}

export function useGetById(id: string) {
  return useQuery<Order>({
    queryKey: ["order", id],
    queryFn: () => fetchWithAuth(`/orders/${id}`),
    enabled: !!id,
  });
}

export function useCreate() {
  const queryClient = useQueryClient();

  return useMutation<Order, Error, Partial<Order>>({
    mutationFn: (data) =>
      fetchWithAuth("/orders", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
    },
  });
}

export function useUpdate() {
  const queryClient = useQueryClient();

  return useMutation<Order, Error, { id: string; data: Partial<Order> }>({
    mutationFn: ({ id, data }) =>
      fetchWithAuth(`/orders/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      queryClient.invalidateQueries({ queryKey: ["order", id] });
    },
  });
}

export function useDeleteById() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: (id) =>
      fetchWithAuth(`/orders/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
    },
  });
}


export function useGetTotalOrdersQuery() {
  return useQuery<number>({
    queryKey: ["orders", "total"],
    queryFn: () => fetchWithAuth("/orders/total"),
  });
}

export function useGetTotalSalesQuery() {
  return useQuery<number>({
    queryKey: ["orders", "total-sales"],
    queryFn: () => fetchWithAuth("/orders/total-sales"),
  });
}

export function useGetTotalSalesByDateQuery() {
  return useQuery<{ _id: string; totalSales: number }[]>({
    queryKey: ["orders", "total-sales-by-date"],
    queryFn: () => fetchWithAuth("/orders/total-sales-by-date"),
  });
}

export function useGetPaypalClientIdQuery() {
  return useQuery<string>({
    queryKey: ["orders", "paypal-client-id"],
    queryFn: () => fetchWithAuth("/config/paypal"),
  });
}

export function usePayOrderMutation() {
  const queryClient = useQueryClient();

  return useMutation<Order, Error, { id: string; data: Partial<Order> }>({
    mutationFn: ({ id, data }) =>
      fetchWithAuth(`/orders/${id}/pay`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      queryClient.invalidateQueries({ queryKey: ["order", id] });
    },
  });
}

export function useDeliverOrderMutation() {
  const queryClient = useQueryClient();

  return useMutation<Order, Error, { id: string; data: Partial<Order> }>({
    mutationFn: ({ id, data }) =>
      fetchWithAuth(`/orders/${id}/deliver`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      queryClient.invalidateQueries({ queryKey: ["order", id] });
    },
  });
}

export function newOrder(): Order {
  return {
    _id: "",
    user: {
      _id: "",
      username: "",
      email: "",
    },
    orderItems: [],
    shippingAddress: {
      address: "",
      city: "",
      postalCode: "",
      country: "",
    },
    paymentMethod: "",
    itemsPrice: 0,
    shippingPrice: 0,
    taxPrice: 0,
    totalPrice: 0,
    isPaid: false,
    isDelivered: false,
    createdAt: "",
    updatedAt: "",
  };
}

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

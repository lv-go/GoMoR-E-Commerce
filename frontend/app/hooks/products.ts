import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchWithAuth } from "../utils/api";
import { type Product } from "../schemas/product";
import { type Page } from "../schemas/api";

export function useGetPage(offset: number = 0, limit: number = 10) {
  return useQuery<Page<Product>>({
    queryKey: ["products", "page", { offset, limit }],
    queryFn: () => fetchWithAuth(`/products?offset=${offset}&limit=${limit}`),
  });
}

export function useGetById(id: string) {
  return useQuery<Product>({
    queryKey: ["product", id],
    queryFn: () => fetchWithAuth(`/products/${id}`),
    enabled: !!id,
  });
}

export function useCreate() {
  const queryClient = useQueryClient();

  return useMutation<Product, Error, Partial<Product>>({
    mutationFn: (data) =>
      fetchWithAuth("/products", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });
}

export function useUpdate() {
  const queryClient = useQueryClient();

  return useMutation<Product, Error, { id: string; data: Partial<Product> }>({
    mutationFn: ({ id, data }) =>
      fetchWithAuth(`/products/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
      queryClient.invalidateQueries({ queryKey: ["product", id] });
    },
  });
}

export function useDeleteById() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: (id) =>
      fetchWithAuth(`/products/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });
}

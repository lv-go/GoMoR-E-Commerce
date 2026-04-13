import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchWithAuth } from "../utils/fetch-with-auth";
import { type Product, type ProductFilters } from "../schemas/product";
import { type Page, type PageRequest } from "../schemas/api";

export function useGetTopProducts() {
  return useGetPage({ offset: 0, limit: 10, sort: "rating", order: "desc" });
}

export function useGetPage(params: PageRequest & ProductFilters) {
  return useQuery<Page<Product>>({
    queryKey: ["products", "page", JSON.stringify(params)],
    queryFn: () => fetchWithAuth(`/products`, params),
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

export function useCreateReview() {
  const queryClient = useQueryClient();

  return useMutation<Product, Error, { productId: string; rating: number; comment: string }>({
    mutationFn: ({ productId, rating, comment }) =>
      fetchWithAuth(`/products/${productId}/reviews`, {
        method: "POST",
        body: JSON.stringify({ rating, comment }),
      }),
    onSuccess: (_, { productId }) => {
      queryClient.invalidateQueries({ queryKey: ["product", productId] });
    },
  });
}

export function useGetBrands() {
  return useQuery<string[]>({
    queryKey: ["products", "brands"],
    queryFn: () => fetchWithAuth(`/products/brands`),
  });
}

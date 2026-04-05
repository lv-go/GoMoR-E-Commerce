import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchWithAuth } from "../utils/api";
import { type Product, type ProductFilters } from "../schemas/product";
import { type Page, type PageRequest } from "../schemas/api";

export function useGetTopProducts() {
  return useGetPage({ offset: 0, limit: 10, sort: "rating", order: "desc" });
}

export function useGetPage({ offset, limit, sort, order, search, category, price }: PageRequest & ProductFilters) {
  const params = new URLSearchParams();
  if (offset) params.set("offset", offset.toString());
  if (limit) params.set("limit", limit.toString());
  if (sort) params.set("sort", sort);
  if (order) params.set("order", order);
  if (search) params.set("search", search);
  if (category) params.set("category", category);
  if (price) params.set("price", price);

  return useQuery<Page<Product>>({
    queryKey: ["products", "page", { offset, limit, sort, order, search, category, price }],
    queryFn: () => fetchWithAuth(`/products?${params.toString()}`),
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

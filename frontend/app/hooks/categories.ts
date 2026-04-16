import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchWithAuth } from "../utils/fetch-with-auth";
import { type Category } from "../schemas/category";
import { type Page, type PageRequest } from "../schemas/api";

export function useGetPage(params: PageRequest = {}) {
  return useQuery<Page<Category>>({
    queryKey: ["categories", "page", JSON.stringify(params)],
    queryFn: () => fetchWithAuth(`/categories`, { params }),
  });
}

export function useGetById(id: string) {
  return useQuery<Category>({
    queryKey: ["category", id],
    queryFn: () => fetchWithAuth(`/categories/${id}`),
    enabled: !!id,
  });
}

export function useCreate() {
  const queryClient = useQueryClient();

  return useMutation<Category, Error, Partial<Category>>({
    mutationFn: (data) =>
      fetchWithAuth("/categories", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["categories"] });
    },
  });
}

export function useUpdate() {
  const queryClient = useQueryClient();

  return useMutation<Category, Error, { id: string; data: Partial<Category> }>({
    mutationFn: ({ id, data }) =>
      fetchWithAuth(`/categories/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["categories"] });
      queryClient.invalidateQueries({ queryKey: ["category", id] });
    },
  });
}

export function useDeleteById() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: (id) =>
      fetchWithAuth(`/categories/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["categories"] });
    },
  });
}

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchWithAuth } from "../utils/api";
import { type User } from "../schemas/user";
import { type Page } from "../schemas/api";

export function useGetPage(offset: number = 0, limit: number = 10) {
  return useQuery<Page<User>>({
    queryKey: ["users", "page", { offset, limit }],
    queryFn: () => fetchWithAuth(`/users?offset=${offset}&limit=${limit}`),
  });
}

export function useGetById(id: string) {
  return useQuery<User>({
    queryKey: ["user", id],
    queryFn: () => fetchWithAuth(`/users/${id}`),
    enabled: !!id,
  });
}

export function useCreate() {
  const queryClient = useQueryClient();

  return useMutation<User, Error, Partial<User>>({
    mutationFn: (data) =>
      fetchWithAuth("/users", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
}

export function useUpdate() {
  const queryClient = useQueryClient();

  return useMutation<User, Error, { id: string; data: Partial<User> }>({
    mutationFn: ({ id, data }) =>
      fetchWithAuth(`/users/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      queryClient.invalidateQueries({ queryKey: ["user", id] });
    },
  });
}

export function useDeleteById() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: (id) =>
      fetchWithAuth(`/users/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
}

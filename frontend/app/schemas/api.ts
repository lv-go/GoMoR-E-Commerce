export interface Page<T> {
  items: T[];
  total: number;
  page: number;
  size: number;
  totalPages: number;
}

export interface PaginationParams {
  offset?: number;
  limit?: number;
}

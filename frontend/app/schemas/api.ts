export interface Page<T> {
  items: T[];
  total: number;
  page: number;
  size: number;
  totalPages: number;
}

export interface PageRequest {
  offset?: number;
  limit?: number;
}

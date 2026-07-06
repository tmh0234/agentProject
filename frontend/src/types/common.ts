export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface PageResult<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}

export interface PageQuery {
  page: number;
  pageSize: number;
  keyword?: string | undefined;
  status?: number | undefined;
}

export type StatusValue = 0 | 1;

export interface SelectOption<T extends string | number = number> {
  label: string;
  value: T;
}

export interface TagView {
  path: string;
  title: string;
  name: string;
}

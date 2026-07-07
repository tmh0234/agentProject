import type { StatusValue } from "./common";

export interface ProductCategory {
  id: number;
  parentId: number;
  name: string;
  status: StatusValue;
  sort: number;
  children?: ProductCategory[];
  createdAt: string;
  updatedAt: string;
}

export interface ProductCategoryForm {
  parentId: number;
  name: string;
  status: StatusValue;
  sort: number;
}

export interface Brand {
  id: number;
  name: string;
  logo: string;
  status: StatusValue;
  sort: number;
  createdAt: string;
  updatedAt: string;
}

export interface BrandForm {
  name: string;
  logo: string;
  status: StatusValue;
  sort: number;
}

export type SkuSpec = Record<string, string>;

export interface ProductSku {
  id: number;
  productId: number;
  code: string;
  specs: SkuSpec;
  priceCents: number;
  stock: number;
  status: StatusValue;
  createdAt: string;
  updatedAt: string;
}

export interface ProductSkuForm {
  id?: number;
  code: string;
  specs: SkuSpec;
  priceCents: number;
  stock: number;
  status: StatusValue;
}

export interface Product {
  id: number;
  name: string;
  code: string;
  categoryId: number;
  brandId: number;
  mainImage: string;
  minPriceCents: number;
  maxPriceCents: number;
  status: StatusValue;
  description: string;
  skus: ProductSku[];
  createdAt: string;
  updatedAt: string;
}

export interface ProductForm {
  name: string;
  code: string;
  categoryId: number;
  brandId: number;
  mainImage: string;
  status: StatusValue;
  description: string;
  skus: ProductSkuForm[];
}

export interface SpecPair {
  key: string;
  value: string;
}

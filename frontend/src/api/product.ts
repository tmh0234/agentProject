import { request, requestPage } from "./http";
import type { PageQuery, PageResult } from "@/types/common";
import type {
  Brand,
  BrandForm,
  Product,
  ProductCategory,
  ProductCategoryForm,
  ProductForm,
  ProductSku,
  ProductSkuForm
} from "@/types/product";

export function getCategoriesApi(params?: Partial<PageQuery>): Promise<ProductCategory[]> {
  return request<ProductCategory[]>({ url: "/product-categories", method: "GET", params });
}

export function createCategoryApi(data: ProductCategoryForm): Promise<ProductCategory> {
  return request<ProductCategory>({ url: "/product-categories", method: "POST", data });
}

export function updateCategoryApi(id: number, data: ProductCategoryForm): Promise<ProductCategory> {
  return request<ProductCategory>({ url: `/product-categories/${id}`, method: "PUT", data });
}

export function deleteCategoryApi(id: number): Promise<void> {
  return request<void>({ url: `/product-categories/${id}`, method: "DELETE" });
}

export function getBrandsApi(params: PageQuery): Promise<PageResult<Brand>> {
  return requestPage<Brand>({ url: "/brands", method: "GET", params });
}

export function createBrandApi(data: BrandForm): Promise<Brand> {
  return request<Brand>({ url: "/brands", method: "POST", data });
}

export function updateBrandApi(id: number, data: BrandForm): Promise<Brand> {
  return request<Brand>({ url: `/brands/${id}`, method: "PUT", data });
}

export function deleteBrandApi(id: number): Promise<void> {
  return request<void>({ url: `/brands/${id}`, method: "DELETE" });
}

export function getProductsApi(params: PageQuery): Promise<PageResult<Product>> {
  return requestPage<Product>({ url: "/products", method: "GET", params });
}

export function getProductApi(id: number): Promise<Product> {
  return request<Product>({ url: `/products/${id}`, method: "GET" });
}

export function createProductApi(data: ProductForm): Promise<Product> {
  return request<Product>({ url: "/products", method: "POST", data });
}

export function updateProductApi(id: number, data: ProductForm): Promise<Product> {
  return request<Product>({ url: `/products/${id}`, method: "PUT", data });
}

export function deleteProductApi(id: number): Promise<void> {
  return request<void>({ url: `/products/${id}`, method: "DELETE" });
}

export function createSkuApi(productId: number, data: ProductSkuForm): Promise<ProductSku> {
  return request<ProductSku>({ url: `/products/${productId}/skus`, method: "POST", data });
}

export function updateSkuApi(productId: number, skuId: number, data: ProductSkuForm): Promise<ProductSku> {
  return request<ProductSku>({ url: `/products/${productId}/skus/${skuId}`, method: "PUT", data });
}

export function deleteSkuApi(productId: number, skuId: number): Promise<void> {
  return request<void>({ url: `/products/${productId}/skus/${skuId}`, method: "DELETE" });
}

export function updateProductStatusApi(id: number, status: number): Promise<void> {
  return request<void>({ url: `/products/${id}/status`, method: "PATCH", data: { status } });
}

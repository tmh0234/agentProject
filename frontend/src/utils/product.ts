import type { StatusValue } from "@/types/common";
import type { Product, ProductForm, ProductSkuForm, SpecPair, SkuSpec } from "@/types/product";

export interface ProductSkuDraft {
  localId: string;
  id?: number;
  code: string;
  specPairs: SpecPair[];
  priceYuan: string;
  stock: number;
  status: StatusValue;
}

export interface ProductDraft {
  name: string;
  code: string;
  categoryId: number | null;
  brandId: number | null;
  mainImage: string;
  status: StatusValue;
  description: string;
  skus: ProductSkuDraft[];
}

export function yuanToCents(value: string | number): number {
  const normalized = typeof value === "number" ? value.toFixed(2) : value.trim();
  if (!/^\d+(\.\d{1,2})?$/.test(normalized)) {
    return 0;
  }

  const [yuan = "0", cents = ""] = normalized.split(".");
  const paddedCents = `${cents}00`.slice(0, 2);
  return Number(yuan) * 100 + Number(paddedCents);
}

export function centsToYuan(value: number): string {
  return (value / 100).toFixed(2);
}

export function specPairsToRecord(pairs: SpecPair[]): SkuSpec {
  return pairs.reduce<SkuSpec>((record, pair) => {
    const key = pair.key.trim();
    const specValue = pair.value.trim();
    if (key && specValue) {
      record[key] = specValue;
    }
    return record;
  }, {});
}

export function specRecordToPairs(record: SkuSpec): SpecPair[] {
  return Object.entries(record).map(([key, value]) => ({ key, value }));
}

export function createEmptyProductDraft(localId: string): ProductDraft {
  return {
    name: "",
    code: "",
    categoryId: null,
    brandId: null,
    mainImage: "",
    status: 0,
    description: "",
    skus: [createEmptySkuDraft(localId)]
  };
}

export function createEmptySkuDraft(localId: string): ProductSkuDraft {
  return {
    localId,
    code: "",
    specPairs: [{ key: "", value: "" }],
    priceYuan: "0.00",
    stock: 0,
    status: 1
  };
}

export function productToDraft(product: Product, createLocalId: () => string = () => "sku"): ProductDraft {
  return {
    name: product.name,
    code: product.code,
    categoryId: product.categoryId,
    brandId: product.brandId,
    mainImage: product.mainImage,
    status: product.status,
    description: product.description,
    skus: product.skus.map<ProductSkuDraft>((sku) => ({
      localId: createLocalId(),
      id: sku.id,
      code: sku.code,
      specPairs: specRecordToPairs(sku.specs),
      priceYuan: centsToYuan(sku.priceCents),
      stock: sku.stock,
      status: sku.status
    }))
  };
}

export function productDraftToPayload(draft: ProductDraft): ProductForm {
  return {
    name: draft.name,
    code: draft.code,
    categoryId: draft.categoryId || 0,
    brandId: draft.brandId || 0,
    mainImage: draft.mainImage,
    status: draft.status,
    description: draft.description,
    skus: draft.skus.map<ProductSkuForm>((sku) => {
      const payload: ProductSkuForm = {
        code: sku.code,
        specs: specPairsToRecord(sku.specPairs),
        priceCents: yuanToCents(sku.priceYuan),
        stock: sku.stock,
        status: sku.status
      };
      if (sku.id) {
        payload.id = sku.id;
      }
      return payload;
    })
  };
}

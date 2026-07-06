import { describe, expect, it } from "vitest";
import {
  centsToYuan,
  productToDraft,
  productDraftToPayload,
  specPairsToRecord,
  specRecordToPairs,
  yuanToCents,
  type ProductDraft
} from "@/utils/product";

describe("product utilities", () => {
  it("converts yuan display value to integer cents", () => {
    expect(yuanToCents("19.90")).toBe(1990);
    expect(yuanToCents(8)).toBe(800);
    expect(yuanToCents("0.01")).toBe(1);
  });

  it("converts integer cents to yuan display value", () => {
    expect(centsToYuan(1990)).toBe("19.90");
    expect(centsToYuan(1)).toBe("0.01");
  });

  it("converts SKU spec pairs to API record and back", () => {
    const record = specPairsToRecord([
      { key: "颜色", value: "黑色" },
      { key: "尺码", value: "L" },
      { key: "", value: "忽略" }
    ]);

    expect(record).toEqual({ 颜色: "黑色", 尺码: "L" });
    expect(specRecordToPairs(record)).toEqual([
      { key: "颜色", value: "黑色" },
      { key: "尺码", value: "L" }
    ]);
  });

  it("maps backend product cents fields to editable yuan draft", () => {
    const draft = productToDraft({
      id: 1,
      name: "测试商品",
      code: "SPU001",
      categoryId: 2,
      brandId: 3,
      mainImage: "",
      minPriceCents: 1299,
      maxPriceCents: 2599,
      status: 1,
      description: "测试",
      skus: [
        {
          id: 10,
          productId: 1,
          code: "SKU001",
          specs: { 颜色: "黑色" },
          priceCents: 1299,
          stock: 8,
          status: 1,
          createdAt: "2026-07-06T10:00:00+08:00",
          updatedAt: "2026-07-06T10:00:00+08:00"
        }
      ],
      createdAt: "2026-07-06T10:00:00+08:00",
      updatedAt: "2026-07-06T10:00:00+08:00"
    });

    expect(draft.code).toBe("SPU001");
    expect(draft.skus[0]?.code).toBe("SKU001");
    expect(draft.skus[0]?.priceYuan).toBe("12.99");
  });

  it("builds product submit payload with backend cents fields", () => {
    const draft: ProductDraft = {
      name: "测试商品",
      code: "SPU001",
      categoryId: 2,
      brandId: 3,
      mainImage: "",
      status: 1,
      description: "测试",
      skus: [
        {
          localId: "sku-1",
          code: "SKU001",
          specPairs: [{ key: "颜色", value: "黑色" }],
          priceYuan: "12.99",
          stock: 8,
          status: 1
        }
      ]
    };

    expect(productDraftToPayload(draft)).toEqual({
      name: "测试商品",
      code: "SPU001",
      categoryId: 2,
      brandId: 3,
      mainImage: "",
      status: 1,
      description: "测试",
      skus: [
        {
          code: "SKU001",
          specs: { 颜色: "黑色" },
          priceCents: 1299,
          stock: 8,
          status: 1
        }
      ]
    });
  });
});

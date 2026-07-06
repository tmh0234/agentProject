import { flushPromises } from "@vue/test-utils";
import { describe, expect, it, vi } from "vitest";

describe("main app render", () => {
  it("renders Element Plus login controls as native form elements", async () => {
    document.body.innerHTML = '<div id="app"></div>';
    globalThis.history = window.history;
    window.history.pushState({}, "", "/");

    await vi.resetModules();
    await import("@/main");
    const { default: router } = await import("@/router");
    await router.isReady();
    for (let index = 0; index < 3; index += 1) {
      await flushPromises();
    }

    const usernameInput = document.querySelector('input[placeholder="请输入用户名"]');
    const passwordInput = document.querySelector('input[placeholder="请输入密码"]');
    const loginButton = Array.from(document.querySelectorAll("button")).find((button) =>
      button.textContent?.includes("登录")
    );

    expect(usernameInput).toBeInstanceOf(HTMLInputElement);
    expect(passwordInput).toBeInstanceOf(HTMLInputElement);
    expect(loginButton).toBeInstanceOf(HTMLButtonElement);
  });
});

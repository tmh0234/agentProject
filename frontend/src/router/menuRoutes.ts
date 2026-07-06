import type { Menu } from "@/types/system";

export function menuRouteName(menu: Menu): string {
  return menu.routeName;
}

export function resolveViewModulePath(componentPath: string): string {
  const normalized = componentPath.endsWith(".vue") ? componentPath : `${componentPath}.vue`;
  const withoutPrefix = normalized.startsWith("views/") ? normalized.slice("views/".length) : normalized;
  return `/src/views/${withoutPrefix}`;
}

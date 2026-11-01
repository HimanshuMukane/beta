
(function () {
  "use strict";

  const components = [];

  function mounted(el) {
    if (!el.isConnected) return false;
    for (let node = el; node; node = node.parentElement) {
      if (node._templPortalOwner && !mounted(node._templPortalOwner)) return false;
    }
    return true;
  }

  function sync(component) {
    for (const el of component.elements) {
      if (mounted(el)) continue;
      component.elements.delete(el);
      component.destroy?.(el);
    }
    for (const el of document.querySelectorAll(component.selector)) {
      if (component.elements.has(el) || !mounted(el)) continue;
      component.elements.add(el);
      component.init?.(el);
    }
  }

  function register(selector, { init, destroy } = {}) {
    const component = { selector, init, destroy, elements: new Set() };
    components.push(component);
    if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", () => sync(component));
    else sync(component);
  }

  new MutationObserver(() => components.forEach(sync))
    .observe(document.documentElement, { childList: true, subtree: true });

  window.templ = window.templ || {};
  window.templ.lifecycle = { register };
})();
